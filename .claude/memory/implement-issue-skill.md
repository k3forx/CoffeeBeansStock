# /implement-issue SKILL

## Overview

The `/implement-issue <number>` SKILL automates GitHub issue implementation with a 3-phase workflow:
- **Phase 1**: Issue analysis (regular mode)
- **Phase 2**: Detailed planning (plan mode)
- **Phase 3**: Implementation with acceptance criteria verification (regular mode)

## Key Features

### 1. Phase Detection
- Automatically detects current mode by checking system reminders for "Plan mode is active"
- Guides user through the workflow: Phase 1 → `/plan` → Phase 2 → approval → Phase 3

### 2. Acceptance Criteria Auto-Check
After implementation completes, AI automatically:
- Extracts acceptance criteria from issue body
- Analyzes implemented code
- Judges each criterion: ✅ Met / ⚠️ Partial / ❌ Not Met
- **Auto-creates PR if all criteria are met (✅)**
- Proposes additional implementation if needed (⚠️/❌)

### 3. End-to-End Automation
- Issue → Branch creation → Implementation → Tests → Acceptance check → Commit → **Auto PR (if all criteria met)**

## Usage Flow

```bash
# Phase 1: Analyze issue
/implement-issue 23

# User enters plan mode
/plan

# Phase 2: Detailed planning
/implement-issue 23
# [AI creates detailed plan]
# [User approves and exits plan mode]

# Phase 3: Implementation
# [AI implements, tests, checks acceptance criteria]
# [AI creates commit]
# [If all criteria ✅: AI automatically creates PR]
# [If any ⚠️/❌: AI asks user for confirmation]
```

## File Locations

- **SKILL definition**: `.claude/skills/implement-issue.md`
- **Permissions**: `.claude/settings.local.json` (updated with git/gh permissions)

## Permissions Added

**Allow:**
- `gh issue comment:*`
- `gh pr view:*`
- `git checkout -b *`
- `git branch *`
- `git status:*`
- `git diff:*`
- `cd frontend && npx tsc*`
- `cd backend && go test*`

**Ask:**
- `gh pr create:*`
- `gh issue close:*`

## Critical Implementation Details

### Acceptance Criteria Extraction
Searches for these section headers in issue body:
- "## 受け入れ条件"
- "## Acceptance Criteria"
- "## 完了条件"
- "## Definition of Done"

Extracts checklist items: `- [ ] ...` or `- [x] ...`

### AI Judgment Logic
For each criterion:
1. Extract keywords (e.g., "登録", "バリデーション")
2. Search implemented files with Grep
3. Read relevant code sections
4. Analyze implementation completeness
5. Assign status with evidence (file:line references)

### Error Handling
- Issue not found → show error, suggest `gh issue list`
- Test failures → offer options (commit anyway / fix / abort)
- Missing acceptance criteria → skip check with warning
- Git conflicts → provide recovery options

## Best Practices

1. **Always use the 3-phase workflow** - don't skip plan mode for complex tasks
2. **Take acceptance criteria seriously** - they define success
3. **Don't ignore test failures** - investigate and fix
4. **Clear commit messages** - include issue number and summary

## Common Pitfalls

- ❌ Running implementation without plan mode (for complex features)
- ❌ Ignoring acceptance criteria check results
- ❌ Committing with test failures
- ❌ Forgetting to reference issue number in commits

## Testing

Test with simple issue first:
```bash
/implement-issue <simple-issue-number>
```

Then try complex multi-file changes to validate full workflow.

## Future Enhancements (Phase 2+)

- Auto-update issues with implementation status
- Metrics collection (implementation time, success rate)
- Template system for common patterns
- Batch implementation for multiple issues
- Advanced AI judgment with test coverage analysis
