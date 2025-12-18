# Pipeline Example

This example demonstrates a small file-processing pipeline using the extracted `pkg/processing`, `pkg/fileops`, and `pkg/resources` packages. It shows how to compose a `processing.Pipeline` with stages, use a `ResourceManager` for temporary resources, and perform safe file operations.

Layout
- `example-code/main.go` — runnable example demonstrating pipeline composition (uses local packages)

Notes
- The example is intentionally small and uses the public interfaces described in the package READMEs.
- This example is documentation-focused and may need project module setup (`go mod`) to run as a standalone program.
