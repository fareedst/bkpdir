# STDD Token Timeline Storyboard

| Step | Highlighted Tokens | Description |
| --- | --- | --- |
| 1 | `[REQ-CFG_005]` | Requirement captured for layered configuration inheritance. |
| 2 | `[ARCH-CFG_005]` | Architecture decision links back to the requirement. |
| 3 | `[IMPL-EXCLUDE_MERGE_FIX]` | Implementation enforces merge behavior. |
| 4 | `[TEST-EXCLUDE_MERGE]`, `config.go::applyMergeStrategies` | Tests and code reference the same tokens. |
| 5 | `[REQ-LIST_LIMIT]` | Parallel requirement for list command limit. |
| 6 | `[ARCH-LIST_LIMIT]` | Architecture decision for list limit. |
| 7 | `[IMPL-LIST_LIMIT]` | Implementation controlling CLI list limit. |
| 8 | `[TEST-LIST_LIMIT]`, `list.go::applyLimit` | Tests/code validate and implement the behavior. |

Use this storyboard to animate token hand-offs; each frame should fade in the highlighted token and show arrows to the next layer.

