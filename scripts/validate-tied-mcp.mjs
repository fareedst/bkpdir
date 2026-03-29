#!/usr/bin/env node
/**
 * Run tied-yaml MCP tools over stdio (same entrypoint as Cursor): yaml_index_validate,
 * then tied_validate_consistency. Exits non-zero on failure.
 *
 * Uses newline-delimited JSON-RPC (same as @modelcontextprotocol/sdk stdio transport).
 *
 * Env:
 *   TIED_YAML_MCP_JS  — path to mcp-server/dist/index.js (required if not in .cursor/mcp.json)
 *   TIED_BASE_PATH    — TIED data root (default: <repo>/tied)
 */
import { spawn } from "node:child_process";
import { existsSync, readFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));

function findRepoRoot(startDir) {
  let d = resolve(startDir);
  for (;;) {
    if (existsSync(join(d, "go.mod"))) return d;
    const p = dirname(d);
    if (p === d) {
      throw new Error("validate-tied-mcp: could not find repo root (go.mod)");
    }
    d = p;
  }
}

function loadMcpConfig(repoRoot) {
  const p = join(repoRoot, ".cursor", "mcp.json");
  if (!existsSync(p)) return null;
  const j = JSON.parse(readFileSync(p, "utf8"));
  const ty =
    j?.mcpServers?.["tied-yaml"] ??
    j?.mcpServers?.["project-0-bkpdir-tied-yaml"];
  if (!ty || ty.type !== "stdio") return null;
  const command = ty.command;
  const args = ty.args ?? [];
  const env = { ...process.env, ...(ty.env ?? {}) };
  return { command, args, env };
}

class NDJSONMCPClient {
  constructor(child) {
    this.child = child;
    this.buf = "";
    this.pending = new Map();
    this.nextId = 1;
    child.stdout.setEncoding("utf8");
    child.stdout.on("data", (chunk) => this._feed(chunk));
    child.stderr.setEncoding("utf8");
    child.stderr.on("data", (ch) => process.stderr.write(ch));
  }

  _feed(chunk) {
    this.buf += chunk;
    for (;;) {
      const i = this.buf.indexOf("\n");
      if (i < 0) return;
      const line = this.buf.slice(0, i).replace(/\r$/, "");
      this.buf = this.buf.slice(i + 1);
      if (!line.trim()) continue;
      let msg;
      try {
        msg = JSON.parse(line);
      } catch {
        continue;
      }
      if (msg.id != null && this.pending.has(msg.id)) {
        const { resolve, reject } = this.pending.get(msg.id);
        this.pending.delete(msg.id);
        if (msg.error) reject(new Error(msg.error.message || JSON.stringify(msg.error)));
        else resolve(msg.result);
      }
    }
  }

  sendLine(obj) {
    this.child.stdin.write(JSON.stringify(obj) + "\n");
  }

  request(method, params) {
    const id = this.nextId++;
    return new Promise((resolve, reject) => {
      this.pending.set(id, { resolve, reject });
      this.sendLine({ jsonrpc: "2.0", id, method, params });
    });
  }

  notify(method, params = {}) {
    this.sendLine({ jsonrpc: "2.0", method, params });
  }
}

function unwrapToolResult(result) {
  if (!result || typeof result !== "object") return result;
  const content = result.content;
  if (Array.isArray(content) && content[0]?.type === "text" && content[0]?.text) {
    try {
      return JSON.parse(content[0].text);
    } catch {
      return result;
    }
  }
  return result;
}

async function main() {
  const repoRoot = findRepoRoot(__dirname);
  const cfg = loadMcpConfig(repoRoot);

  const serverJs = process.env.TIED_YAML_MCP_JS || (cfg?.args?.[0] ?? "");
  if (!serverJs || !existsSync(serverJs)) {
    console.error(
      "validate-tied-mcp: Set TIED_YAML_MCP_JS to mcp-server/dist/index.js, or configure .cursor/mcp.json tied-yaml stdio args[0].",
    );
    console.error("Expected file missing:", serverJs || "(empty)");
    process.exit(2);
  }

  let tiedBase =
    process.env.TIED_BASE_PATH ||
    cfg?.env?.TIED_BASE_PATH ||
    join(repoRoot, "tied");
  if (!resolve(tiedBase).startsWith("/")) {
    tiedBase = resolve(repoRoot, tiedBase);
  }

  const command = cfg?.command || "node";
  const args = [serverJs];
  const childEnv = { ...process.env, ...cfg?.env, TIED_BASE_PATH: tiedBase };

  const child = spawn(command, args, {
    env: childEnv,
    stdio: ["pipe", "pipe", "pipe"],
  });

  const client = new NDJSONMCPClient(child);

  const onExit = new Promise((_, rej) =>
    child.on("exit", (code, sig) => {
      if (code !== 0 && code !== null)
        rej(new Error(`MCP server exited (code=${code} signal=${sig})`));
    }),
  );

  try {
    await Promise.race([
      (async () => {
        await client.request("initialize", {
          protocolVersion: "2024-11-05",
          capabilities: {},
          clientInfo: { name: "validate-tied-mcp", version: "1.0.0" },
        });
        client.notify("notifications/initialized", {});

        const idxRaw = await client.request("tools/call", {
          name: "yaml_index_validate",
          arguments: {},
        });
        const idx = unwrapToolResult(idxRaw);
        console.log("yaml_index_validate:", JSON.stringify(idx, null, 2));

        let idxFail = false;
        if (idx && typeof idx === "object") {
          for (const v of Object.values(idx)) {
            if (v && typeof v === "object" && v.valid === false) idxFail = true;
          }
        }
        if (idxFail) {
          console.error("validate-tied-mcp: yaml_index_validate reported invalid index.");
          process.exitCode = 1;
        }

        const consRaw = await client.request("tools/call", {
          name: "tied_validate_consistency",
          arguments: {
            include_detail_files: true,
            include_pseudocode: true,
            require_detail_record: true,
          },
        });
        const cons = unwrapToolResult(consRaw);
        const ok = cons?.ok === true;
        console.log("tied_validate_consistency ok:", ok);
        if (!ok) {
          console.error(JSON.stringify(cons, null, 2));
          process.exitCode = 1;
        }
      })(),
      onExit,
    ]);
  } catch (e) {
    console.error("validate-tied-mcp:", e.message || e);
    process.exitCode = 1;
  } finally {
    try {
      child.stdin.end();
    } catch {
      /* ignore */
    }
    child.kill("SIGTERM");
  }
}

main().then(() => process.exit(process.exitCode || 0));
