# STDD Trace Data Schema (Token-First)

This schema defines the token-first data model for visualizations that show `[REQ:*] → [ARCH:*] → [IMPL:*] → TEST/CODE` chains.

## Schema

```json
{
  "schemaVersion": "0.1",
  "generatedAt": "ISO-8601 timestamp",
  "sourceDocs": ["tied/semantic-tokens.yaml", "tied/semantic-tokens.md", "..."],
  "nodes": [
    {
      "id": "string",                // unique node id
      "token": "REQ:CFG_005",        // semantic token
      "type": "REQ|ARCH|IMPL|TEST|CODE",
      "title": "Human-friendly name",
      "status": "planned|in_progress|complete",
      "source": "path or doc anchor",
      "description": "short optional text",
      "styling": {
        "lane": "REQ|ARCH|IMPL|TEST|CODE",
        "colorHint": "#RRGGBB or semantic name",
        "weight": 1                  // relative prominence
      },
      "examples": [
        {"path": "config.go", "symbol": "applyMergeStrategies"}
      ]
    }
  ],
  "edges": [
    {
      "from": "node-id",
      "to": "node-id",
      "relation": "satisfies|implements|validated_by|tested_by|located_in",
      "note": "optional"
    }
  ]
}
```

## Token-First Principles
- Tokens are the primary visual elements; styling hints keep tokens prominent.
- Nodes carry lane and color hints so layered and timeline visuals can render tokens consistently.
- Edges capture traceability; visuals should label edges lightly to keep focus on tokens.

## Source Queries (reference)
- Requirements/architecture/implementation: use `tied/requirements.yaml`, `tied/architecture-decisions.yaml`, `tied/implementation-decisions.yaml` and detail YAML under `tied/`; human digests in `tied/requirements.md`, `tied/architecture-decisions.md`, `tied/implementation-decisions.md`.
- Token registry: `tied/semantic-tokens.yaml` and `tied/semantic-tokens.md` for canonical token list and types.
- Code/tests: grep for `\[REQ:` `\[ARCH:` `\[IMPL:` in `*.go`, `*_test.go` to locate anchors for CODE/TEST nodes and examples.
- Sample chain selection: choose representative chains, then resolve doc anchors and code symbols for each hop.

## Consistency Checks (lightweight)
- Node/edge count sanity: nodes > 0; edges connect valid node ids; lane/styling present for all nodes.
- Sample chains: every hop includes a token and valid node id; chain length ≥ 3; ends in TEST or CODE.
- Registry alignment (sampled): confirm sampled tokens exist in `tied/semantic-tokens.yaml`; note that sampled graphs are not full coverage unless explicitly generated.

## Validation Checklist
- Schema version present.
- Node/edge counts align with token registry samples.
- Each sample chain includes tokens at every hop (REQ→ARCH→IMPL→TEST/CODE).
- Styling hints exist for all nodes to support token prominence in visuals.
- Source queries documented for TIED YAML / digests and code/test anchors, with regex for `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`.
- Consistency checks run: token count vs registry; presence of required fields; sample chains validated.

