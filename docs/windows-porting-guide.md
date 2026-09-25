# Windows portability changes — porting and merge guide

This document captures **all** changes introduced to make the `bkpdir` TIED client build and run on Windows, as recorded in a single commit on this repository. Use it to reproduce the same result on another machine or to **merge** these changes onto a branch that shares the same starting point but has diverged with other work.

## Commit anchors

| Role | SHA | Message |
| --- | --- | --- |
| **Merge base (parent)** | `7ec0cec7700ea560b46c6e649973bd39eace59dd` | `refactor(cli): share production root command with tests` |
| **Windows portability commit** | `db66a3b7ed774b4a46f55b6472982d48c766fec3` | `win11` |

**Diff scope:** one commit, parent → `db66a3b7` (six paths touched; see [Changed files](#changed-files)).

**TIED YAML:** this commit does **not** modify `tied/` indexes or detail files. Porting is documentation + Go CLI + one composition test only.

**Binary artifact:** the commit adds `bkpdir.exe` (~9.5 MB). That file is a **local build output**. On another clone, prefer `go build -o bkpdir.exe .` instead of copying the binary. Omit `bkpdir.exe` from git when merging unless you explicitly want a checked-in Windows binary.

---

## What this change set accomplishes

1. **Document** that macOS/Linux remain primary CI targets while Windows can build and run from source with known test limitations.
2. **Document** PowerShell workflows (`go mod download`, `go build`, `go install`, smoke tests, Git Bash on `PATH`) instead of assuming `make`.
3. **Fix CLI usability on Windows/Cobra:** bare `bkpdir --list` (directory archives, backward compatibility) is rewritten to the `list` subcommand before Cobra parsing, so the user is not forced to pass a string value to the `--list` flag.
4. **Lock behavior** with a composition test for bare `--list`.
5. **Clarify** user-facing platform compatibility in `docs/user/specification.md`.
6. **Normalize ZIP entry paths at write time** so archive member names always use forward slashes (`/`), matching PKZIP convention and macOS/Linux `unzip` / bkpdir directory snapshots (`OracleZipEntryPath` in `addFileToZip*`).

### ZIP entry paths (cross-platform archives)

Directory archives store each file’s path in the ZIP central directory. On Windows, `filepath.Rel` historically produced backslashes in those names, which Mac tools treat as literal characters instead of directories.

After the write-path fix in `archive.go`:

- **New archives** created on any OS use `/` in member names (same as Mac/Linux today).
- **Archives created on Windows before the fix** may still contain `\` in names; re-create the archive on a fixed build or use an offline repair utility if you maintain one separately.

No change to archive **filenames** on disk (`bkpdir-YYYY-MM-DD-hh-mm.zip` grammar is unchanged).

---

## Recommended merge strategies

Choose based on where your target branch points.

### A. Target branch is exactly at the parent commit

If `HEAD` is `7ec0cec7700ea560b46c6e649973bd39eace59dd` (or you can reset/rebase your Windows work to start there):

```bash
git cherry-pick db66a3b7ed774b4a46f55b6472982d48c766fec3
```

Drop the binary from the index if you do not want it tracked:

```bash
git reset HEAD bkpdir.exe 2>/dev/null || true
git restore --staged bkpdir.exe 2>/dev/null || true
echo "bkpdir.exe" >> .gitignore   # only if not already ignored
```

### B. Target branch diverged from the same parent (parallel work)

Both your branch and the Windows commit should share ancestor **`7ec0cec`**. Merge the Windows commit on top of your work:

```bash
# On your feature branch (contains other changes after 7ec0cec)
git fetch origin   # if the Windows commit lives on a remote
git cherry-pick db66a3b7ed774b4a46f55b6472982d48c766fec3
```

Resolve conflicts file-by-file using [Manual application checklist](#manual-application-checklist) below, then:

```bash
git cherry-pick --continue
```

If cherry-pick is awkward (many conflicts), apply a patch instead:

```bash
git format-patch -1 db66a3b7ed774b4a46f55b6472982d48c766fec3 --stdout > windows-portability.patch
# On target branch at or after 7ec0cec:
git apply --3way windows-portability.patch
# Edit conflicts, then commit with message referencing db66a3b7
```

### C. Target repo has no access to this SHA

Export from a clone that contains `db66a3b7`:

```bash
git diff 7ec0cec7700ea560b46c6e649973bd39eace59dd db66a3b7ed774b4a46f55b6472982d48c766fec3 \
  -- DEVELOPMENT.md README.md docs/user/specification.md main.go main_cli_composition_test.go \
  > windows-portability-no-exe.diff
```

Apply on the other system at the matching parent, then commit.

---

## Changed files

| File | Action | Purpose |
| --- | --- | --- |
| `DEVELOPMENT.md` | Edit | Prerequisites (Go 1.21+, optional Make), Windows quick start, troubleshooting |
| `README.md` | Edit | Split Installation: macOS/Linux vs Windows; `go run .`; dev note |
| `docs/user/specification.md` | Edit | Platform compatibility paragraph for Windows source builds |
| `main.go` | Edit | `normalizeBareListFlag`; call from `executeWithAutoDetection`; help/flag text |
| `main_cli_composition_test.go` | Edit | `TestComposition_RootBareListFlagListsArchives` |
| `bkpdir.exe` | Added (optional) | Local build; skip when porting |

---

## Manual application checklist

Apply in this order if not using a clean cherry-pick.

### 1. `DEVELOPMENT.md`

- **Prerequisites:** Go **1.21+** (reference `go.mod`); Make **optional on Windows**; add paragraph that primary targets are macOS/Linux and Windows can build/run but root-module tests are not Windows-clean.
- **New subsection `### Windows`** under building: PowerShell `go mod download`, `go build -o bkpdir.exe .`, `go install .`, smoke tests in `pkg/fileops` and `pkg/formatter`, prepend Git to `PATH` for bash-using tests, note on `go mod verify` vs local `replace` modules.
- **Troubleshooting:** Go **1.21+**; three Windows bullets (no `make`, `-ldflags`/PowerShell quoting, `/tmp` test failures).

### 2. `README.md`

- **Installation:** Intro sentence on primary platforms vs Windows; split **macOS / Linux** (keep `make`; change `go run main.go` → `go run .`) and **Windows (PowerShell)** blocks; `go install .` + PATH note; `go run . --help` without install.
- **Development:** One sentence pointing Windows users to Installation instead of `make` unless Make + Git Bash/WSL.

### 3. `docs/user/specification.md`

- In **Platform compatibility**, append: Windows users may build from source (README); automated tests validated primarily on macOS and Linux.

### 4. `main.go`

**Help text** (root command long description): add line documenting `bkpdir --list` for directory archives (backward compatibility).

**Comment** above `--list` handling in `Run`: note bare `--list` is normalized separately.

**Flag** `--list` description: clarify file-specific list vs `list` / bare `--list` for archives.

**New function** `normalizeBareListFlag(args []string) []string` with tokens `[IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY]`:

- Walk `args`; when token is exactly `--list`:
  - If next arg exists and does **not** start with `-`, treat as `--list <file>` (keep both).
  - Otherwise replace `--list` with subcommand name `list`.
- Return new slice.

**`executeWithAutoDetection`:** first line after function entry: `args = normalizeBareListFlag(args)`.

Reference implementation (merge into existing `main.go`; adjust only if your branch moved adjacent code):

```go
func normalizeBareListFlag(args []string) []string {
	if len(args) == 0 {
		return args
	}
	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--list" {
			hasFileArg := i+1 < len(args) && !strings.HasPrefix(args[i+1], "-")
			if hasFileArg {
				out = append(out, arg, args[i+1])
				i++
				continue
			}
			out = append(out, "list")
			continue
		}
		out = append(out, arg)
	}
	return out
}
```

Ensure `strings` is imported (already required elsewhere in `main.go`).

### 5. `main_cli_composition_test.go`

Add test **`TestComposition_RootBareListFlagListsArchives`** after the root persistent config flag test:

- Temp dir with `.bkpdir.yml` (`archive_dir_path: ./archives`) and empty `archives/` directory.
- `Chdir` into temp dir; defer restore.
- `newRootCommand()`, capture stdout/stderr, `executeWithAutoDetection(root, []string{"--list"})` must succeed (empty archive list is OK).

---

## Verification on the target system

### Windows (PowerShell)

```powershell
cd path\to\bkpdir
go mod download
go build -o bkpdir.exe .
.\bkpdir.exe --help
.\bkpdir.exe --list          # must list directory archives, not fail on missing flag value
go test -run TestComposition_RootBareListFlagListsArchives .
cd pkg\fileops; go test ./...
cd ..\formatter; go test ./...
```

Optional full suite (expect failures unrelated to this port):

```powershell
go test ./...
```

### macOS / Linux (regression)

```bash
go test -run TestComposition_RootBareListFlagListsArchives .
go test ./...
make test   # if Makefile workflow unchanged
bkpdir --list # or go run . --list
```

---

## Conflict hints when merging with other CLI work

| Area | Likely overlap | Resolution principle |
| --- | --- | --- |
| `main.go` | `newRootCommand`, `executeWithAutoDetection`, flag definitions | Keep **both** your branch logic and **`normalizeBareListFlag`** call at start of `executeWithAutoDetection` |
| `main_cli_composition_test.go` | Other composition tests | Add new test function; do not remove existing tests |
| `README.md` / `DEVELOPMENT.md` | Installation or dev docs edited elsewhere | Merge prose; preserve Windows PowerShell blocks and primary-platform disclaimer |
| `bkpdir.exe` | N/A | Prefer not to commit; build locally |

After merge, run the composition test and a manual `bkpdir --list` on Windows if available.

---

## Source of truth in git

To inspect the canonical diff at any time:

```bash
git show db66a3b7ed774b4a46f55b6472982d48c766fec3
git diff 7ec0cec7700ea560b46c6e649973bd39eace59dd db66a3b7ed774b4a46f55b6472982d48c766fec3
```

**Document version:** derived from commit `db66a3b7ed774b4a46f55b6472982d48c766fec3` (2026-09-21).
