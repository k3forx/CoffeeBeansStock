# CoffeeBeansStock Project Memory

## Custom SKILLs

### /implement-issue
Automates GitHub issue implementation with 3-phase workflow:
1. Issue analysis (regular mode)
2. Detailed planning (plan mode)
3. Implementation with AI-powered acceptance criteria verification

**Key features**:
- AI automatically checks if all acceptance criteria are met (✅/⚠️/❌)
- **Auto-creates PR when all criteria are satisfied (all ✅)**
- Proposes additional implementation if criteria are not met (⚠️/❌)

See: `memory/implement-issue-skill.md` for detailed documentation.

### /create-issue
Creates GitHub issues interactively with title, body, labels, and assignees.

## Project Structure

- **Frontend**: React Native with Expo (TypeScript)
- **Backend**: Go with sqlc, PostgreSQL
- **Issue Templates**: `.github/ISSUE_TEMPLATE/feature.yml` with structured acceptance criteria

## Development Patterns

### Issue-Driven Development
Issues use detailed templates with:
- Phase (MVP/Phase 2/Phase 3)
- Implementation area (Frontend/Backend/Both)
- Priority (High/Medium/Low)
- Acceptance criteria (checklist format)
- Technical details

The `/implement-issue` SKILL automates implementation of these structured issues.

## Git Workflow

Main branch: `main`
Feature branches: `feature/issue-<number>-<description>`

## Permissions Configuration

Managed in `.claude/settings.local.json`:
- Git operations (branch, checkout, add, commit) - mostly "allow" or "ask"
- GitHub CLI (gh issue, gh pr) - "allow" for read, "ask" for write
- Tests (npx tsc, go test) - "allow"
- Destructive operations - "deny"
