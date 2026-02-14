---
name: create-issue
description: Create a GitHub issue interactively with title, body, labels, and assignees
---

# GitHub Issue Creation Skill

You are helping the user create a GitHub issue for this repository.

## Steps

1. **Gather information** - Ask the user for the following in a single, friendly question:
   - Issue title (required)
   - Issue description/body (required)
   - Labels (optional) - comma-separated
   - Assignees (optional) - comma-separated or GitHub usernames
   - Milestone (optional)

   Present this as a structured question using AskUserQuestion tool with multiple questions.

2. **Create the issue** - Once you have the required information (at minimum title and body), use the `gh issue create` command:

   ```bash
   gh issue create --title "TITLE" --body "$(cat <<'EOF'
   BODY
   EOF
   )" [OPTIONS]
   ```

3. **Return the result** - Show the user:
   - The created issue URL
   - A confirmation message
   - Option to open in browser if desired

## GitHub CLI Commands

Primary command:
```bash
gh issue create [flags]
```

Available flags:
- `-t, --title <string>`: Issue title (required)
- `-b, --body <string>`: Issue body/description (required)
- `-l, --label <name>`: Add label (can be used multiple times)
- `-a, --assignee <login>`: Assign person (can be used multiple times)
- `-m, --milestone <name>`: Add to milestone
- `-w, --web`: Open in web browser after creation

## Tips

- For multi-line bodies, always use heredoc syntax to properly handle newlines and special characters
- Validate that title is not empty before creating
- If user provides labels/assignees as comma-separated, split them into multiple flags
- After creation, offer to open the issue in the browser if the user wants

## Example

```bash
gh issue create \
  --title "Add user authentication" \
  --body "$(cat <<'EOF'
We need to implement user authentication with the following features:
- Login/logout
- Password reset
- Session management
EOF
)" \
  --label "enhancement" \
  --label "priority-high" \
  --assignee "username"
```
