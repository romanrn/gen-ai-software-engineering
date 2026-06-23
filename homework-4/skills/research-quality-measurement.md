# Skill: Research Quality Measurement

Use this skill to assign a quality level to codebase research before passing it to the Bug Planner.

## Quality Levels

### GOLD
- All file:line references verified against source
- All code snippets match source exactly
- Root cause clearly identified
- No discrepancies found

### SILVER
- ≥90% of file:line references verified
- Minor snippet differences (whitespace, formatting)
- Root cause identified with minor gaps
- Discrepancies documented

### BRONZE
- ≥70% of file:line references verified
- Some snippets differ from source
- Root cause partially identified
- Discrepancies documented, Bug Planner must re-verify

### FAIL
- <70% of file:line references verified
- Critical snippets missing or wrong
- Root cause not identified
- Research must be redone before proceeding

## How to Apply

1. Read every file:line reference in `codebase-research.md` and open the actual source file at that line
2. Compare the snippet in the research against the actual code
3. Count verified vs total references
4. Assign the level that matches the results
5. Document each discrepancy with: reference claimed → actual value found

## Output Format

In `verified-research.md`, include:

```
Research Quality: <GOLD|SILVER|BRONZE|FAIL>
Verified: X/Y references
```