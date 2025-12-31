# STDD Trace Data Schema (Token-First)

This schema defines the token-first data model for visualizations that show `[REQ:*] → [ARCH:*] → [IMPL:*] → TEST/CODE` chains.

## Schema

```json
{
  "schemaVersion": "0.1",
  "generatedAt": "ISO-8601 timestamp",
  "sourceDocs": ["stdd/semantic-tokens.md", "..."],
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
- Requirements/architecture/implementation: parse `stdd/requirements.md`, `stdd/architecture-decisions.md`, `stdd/implementation-decisions.md` for `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`.
- Token registry: `stdd/semantic-tokens.md` for canonical token list and types.
- Code/tests: grep for `\[REQ:` `\[ARCH:` `\[IMPL:` in `*.go`, `*_test.go` to locate anchors for CODE/TEST nodes and examples.
- Sample chain selection: choose representative chains, then resolve doc anchors and code symbols for each hop.

## Consistency Checks (lightweight)
- Node/edge count sanity: nodes > 0; edges connect valid node ids; lane/styling present for all nodes.
- Sample chains: every hop includes a token and valid node id; chain length ≥ 3; ends in TEST or CODE.
- Registry alignment (sampled): confirm sampled tokens exist in `stdd/semantic-tokens.md`; note that sampled graphs are not full coverage unless explicitly generated.

## Validation Checklist
- Schema version present.
- Node/edge counts align with token registry samples.
- Each sample chain includes tokens at every hop (REQ→ARCH→IMPL→TEST/CODE).
- Styling hints exist for all nodes to support token prominence in visuals.
- Source queries documented for STDD docs and code/test anchors, with regex for `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`.
- Consistency checks run: token count vs registry; presence of required fields; sample chains validated.

