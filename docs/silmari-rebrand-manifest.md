# Silmari Rebrand Manifest

## Summary
This manifest bootstraps the validation contract for the Silmari rebrand.
AF-RB-02 establishes the required manifest shape, scanner CLI, and verification
evidence while later AF-RB slices populate repo-wide audited files and
preserved compatibility identifiers.

## Audited Files
| Path | Action | Verification |
|---|---|---|

## Preserved Non-Silmari Identifiers
| Path | Identifier | Category | Reason |
|---|---|---|---|

## CodeCleanup Passes
| Pass | Result | Notes |
|---|---|---|
| Brand surface pass | not-applicable | Later slices update visible prose, labels, and assets. |
| Compatibility pass | pass | The scanner enforces allowed preservation categories and span-based coverage rules. |
| Mirror and generated-template pass | pass | The scanner checks embedded skill mirror parity without mutating files. |
| Link and asset pass | not-applicable | Link-label and asset relabeling land in later rebrand slices. |
| Formatting and lint pass | pass | The validation slice keeps the manifest, scanner, and tests structured for automation. |
| Verification pass | pass | Scanner contract tests cover the validation behaviors introduced by AF-RB-02. |

## Verification Commands
| Command | Working Directory | Exit Code | Result |
|---|---|---|---|
| python3 -m pytest tests/test_check_silmari_rebrand.py | . | 0 | Passed 10 scanner contract tests for manifest validation, preserved coverage, manifest exclusion, and mirror drift. |
| python3 -m pytest tests/test_collect_silmari_rebrand_inventory.py | . | 0 | Passed the existing AF-RB-01 inventory regression tests after adding the scanner slice. |
| ./scripts/sync-embedded-skills.sh --check | . | 0 | Embedded skills are in sync with the source skill tree. |
| ./scripts/check-silmari-rebrand.sh | . | 1 | Expected failure at this slice because repo-wide audited-file coverage is still incomplete. |

## Deferred Or Excluded
- Repo-wide Silmari copy edits land in later AF-RB slices.
- Full audited-file and preserved-identifier coverage will expand as each
  rebrand surface is reconciled.
