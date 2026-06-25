#!/bin/bash
# Pipeline runner: executes all 6 agents in order for a given bug directory.
# Usage: ./run-pipeline.sh <bug-dir>
#   e.g. ./run-pipeline.sh context/bugs/bug001
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BUG_DIR_ARG="${1:?Usage: ./run-pipeline.sh <bug-dir>}"

# Resolve BUG_DIR to absolute path
if [[ "$BUG_DIR_ARG" = /* ]]; then
  BUG_DIR="$BUG_DIR_ARG"
else
  BUG_DIR="$SCRIPT_DIR/$BUG_DIR_ARG"
fi

[[ -d "$BUG_DIR" ]] || { echo "ERROR: not a directory: $BUG_DIR" >&2; exit 1; }
[[ -f "$BUG_DIR/bug-context.md" ]] || { echo "ERROR: missing $BUG_DIR/bug-context.md" >&2; exit 1; }

ISSUE_ID="$(basename "$BUG_DIR")"
mkdir -p "$BUG_DIR/research"

# ── helpers ────────────────────────────────────────────────────────────────────

divider() {
  echo ""
  echo "══════════════════════════════════════════════════════════════"
  echo "  $*"
  echo "══════════════════════════════════════════════════════════════"
}

# Strip YAML frontmatter (between first and second ---) and substitute BUG_DIR.
build_prompt() {
  local agent_file="$1"
  awk 'BEGIN{f=0} /^---$/{f++; next} f>=2{print}' "$agent_file" \
    | sed "s|\\\$BUG_DIR|$BUG_DIR|g"
}

# Extract model: field from frontmatter, default to claude-sonnet-4-6.
get_model() {
  local model
  model=$(grep '^model:' "$1" | awk '{print $2}' | tr -d '"' | head -1)
  echo "${model:-claude-sonnet-4-6}"
}

# ── run one agent ───────────────────────────────────────────────────────────────

run_agent() {
  local step="$1"
  local name="$2"
  local agent_rel="$3"
  local agent_file="$SCRIPT_DIR/$agent_rel"

  divider "[$step/6] $name"

  local model
  model=$(get_model "$agent_file")
  echo "  model : $model"
  echo "  agent : $agent_rel"
  echo ""

  local prompt
  prompt=$(build_prompt "$agent_file")

  cd "$SCRIPT_DIR"
  claude -p "Execute your task." \
    --model "$model" \
    --system-prompt "$prompt" \
    --dangerously-skip-permissions \
    --add-dir "$SCRIPT_DIR" \
    2>&1

  echo ""
  echo "  ✓ [$step/6] $name — done"
}

# ── pipeline ────────────────────────────────────────────────────────────────────

divider "Pipeline start — $ISSUE_ID"
echo "  BUG_DIR : $BUG_DIR"
echo ""
date

run_agent 1 "Bug Researcher"      "agents/bug-researcher.agent.md"
run_agent 2 "Research Verifier"   "agents/research-verifier.agent.md"
run_agent 3 "Bug Planner"         "agents/bug-planner.agent.md"
run_agent 4 "Bug Fixer"           "agents/bug-fixer.agent.md"
run_agent 5 "Security Verifier"   "agents/security-verifier.agent.md"
run_agent 6 "Unit Test Generator" "agents/unit-test-generator.agent.md"

divider "Pipeline complete — $ISSUE_ID"
date
echo ""
echo "Artifacts:"
find "$BUG_DIR" -name "*.md" | sort | sed 's|^|  |'