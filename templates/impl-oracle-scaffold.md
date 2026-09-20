# Oracle scaffold template (L3 logic verification)

Copy when adding `internal/specmodel/{domain}.go` + `test/specconformance/{domain}_conformance_test.go`.

## 1. Register in `tied/spec/oracle-registry.yaml`

```yaml
  - impl: IMPL-TOKEN
    specmodel_file: internal/specmodel/{domain}.go
    conformance_test: test/specconformance/{domain}_conformance_test.go
    pkg_paths:
      - pkg/{package}
    status: verified
```

## 2. Oracle function (pure semantics from sidecar PRE/POST/BRANCH/ERROR)

```go
// [IMPL-TOKEN] [ARCH-*] [REQ-*]
func OracleExample(state Input) (Output, error) {
  // PRE from sidecar
  // STEP transitions
  // POST from sidecar
}
```

## 3. Conformance test (oracle vs pkg)

```go
func TestExample_OracleMatchesPkg(t *testing.T) {
  // table or gopter property
  // compare specmodel.OracleX(...) with pkg.X(...)
}
```

## 4. Tracking

Run `python3 scripts/seed_logic_verification_tracking.py` after registry update.
