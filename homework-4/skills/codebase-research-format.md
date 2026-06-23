# Skill: Codebase Research Format

Use this skill to structure research output into a consistent, verifiable document.

## Output File

Write to: `$BUG_DIR/research/codebase-research.md`

## Rules

- Every claim about code must be backed by a file:line reference
- Code snippets must be copied verbatim from source — no paraphrasing, no reformatting
- Do not suggest fixes — root cause analysis only
- If the bug cannot be located, state that clearly rather than guessing

## Output Template

Copy this template exactly:

```markdown
# Codebase Research: <issue-id>

## Summary
<one paragraph: observed symptom and suspected root cause>

## Root Cause
<precise technical explanation of why the bug occurs>

## File References
- `src/main.go:NN` — <description of what is on this line>

## Code Snippets
<verbatim code copied from source, no modifications>

## Impact
<what goes wrong at runtime and who is affected>

## References
- <list of every file read during research>
```