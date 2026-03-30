Use GitHub CLI on macOS. The current syntax is `gh release upload <tag> <files>...`, and you can target another repo with `--repo OWNER/REPO`. Existing asset names can be replaced with `--clobber`. ([GitHub CLI][1])

Install and authenticate first:

```bash
brew install gh
```
```bash
which gh
```
```bash
gh auth login
```

Set the repo, tag, and five binaries:

```ux
init: fareedst
name: OWNER
prompt: "Enter the GitHub owner or organization"
```

```ux
allow:
- bkpdir
- github-markdown-test
name: REPO
prompt: "Enter the GitHub repository name"
```

```ux
init: bkpdir-v1.7.2
name: TAG
prompt: "Enter the release tag"
```

```ux
init: bin/bkpdir-macos-amd64
name: FILE1
prompt: "Enter binary file 1"
```

```ux
init: bin/bkpdir-macos-arm64
name: FILE2
prompt: "Enter binary file 2"
```

```ux
init: bin/bkpdir-ubuntu20.04
name: FILE3
prompt: "Enter binary file 3"
```

```ux
init: bin/bkpdir-ubuntu22.04
name: FILE4
prompt: "Enter binary file 4"
```

```ux
init: bin/bkpdir-ubuntu24.04
name: FILE5
prompt: "Enter binary file 5"
```

Then run this on macOS Terminal:

```bash
#!/usr/bin/env bash
set -euo pipefail

OWNER="${OWNER:-your-org}"
REPO="${REPO:-your-repo}"
TAG="${TAG:-v1.7.2}"

FILES=(
  "${FILE1:-./dist/app-darwin-arm64}"
  "${FILE2:-./dist/app-darwin-amd64}"
  "${FILE3:-./dist/app-linux-arm64}"
  "${FILE4:-./dist/app-linux-amd64}"
  "${FILE5:-./dist/app-windows-amd64.exe}"
)

for f in "${FILES[@]}"; do
  [[ -f "$f" ]] || { echo "Missing file: $f" >&2; exit 1; }
done

set -x

if ! gh release view "$TAG" --repo "$OWNER/$REPO" >/dev/null 2>&1; then
  gh release create "$TAG" \
    --repo "$OWNER/$REPO" \
    --title "$TAG" \
    --notes ""
fi

gh release upload "$TAG" \
  "${FILES[@]}" \
  --repo "$OWNER/$REPO" \
  --clobber

echo "Uploaded ${#FILES[@]} binaries to $OWNER/$REPO release $TAG"

set +x
```

That script checks whether the release for the tag exists, creates it if needed, then uploads all five files. `gh release view` can inspect a tag-specific release, and `gh release create` is the CLI command for creating one. ([GitHub CLI][2])

To see recent releases before choosing a tag:

```bash
gh release list --repo "$OWNER/$REPO" --limit 10
```

`gh release list` is the current CLI command for listing releases. ([GitHub CLI][3])

With asset labels shown in the GitHub UI, append `#label` to each filename when uploading; GitHub CLI supports that format on `gh release upload`. ([GitHub CLI][1])

[1]: https://cli.github.com/manual/gh_release_upload?utm_source=chatgpt.com "gh release upload | GitHub CLI"
[2]: https://cli.github.com/manual/gh_release_view?utm_source=chatgpt.com "gh release view - GitHub CLI"
[3]: https://cli.github.com/manual/gh_release?utm_source=chatgpt.com "gh release - GitHub CLI"
