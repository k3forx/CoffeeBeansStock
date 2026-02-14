---
name: create-issue
shortDescription: タイトル、本文、ラベル、担当者を指定してGitHub issueをインタラクティブに作成する
type: user-invocable
---

# GitHub Issue作成スキル

このリポジトリのGitHub issueを作成するお手伝いをします。

## ステップ

1. **情報収集** - 以下の情報をユーザーに尋ねます：
   - Issueタイトル（必須）
   - Issue本文/説明（必須）
   - ラベル（任意） - カンマ区切り
   - 担当者（任意） - カンマ区切りまたはGitHubユーザー名
   - マイルストーン（任意）

   AskUserQuestionツールを使用して、構造化された質問として提示します。

2. **Issueの作成** - 必要な情報（最低限タイトルと本文）が揃ったら、`gh issue create`コマンドを使用します：

   ```bash
   gh issue create --title "タイトル" --body "$(cat <<'EOF'
   本文
   EOF
   )" [オプション]
   ```

3. **結果の返却** - ユーザーに以下を表示します：
   - 作成されたIssueのURL
   - 確認メッセージ
   - 必要に応じてブラウザで開くオプション

## GitHub CLIコマンド

主要コマンド：
```bash
gh issue create [flags]
```

利用可能なフラグ：
- `-t, --title <string>`: Issueタイトル（必須）
- `-b, --body <string>`: Issue本文/説明（必須）
- `-l, --label <name>`: ラベルを追加（複数回使用可能）
- `-a, --assignee <login>`: 担当者を割り当て（複数回使用可能）
- `-m, --milestone <name>`: マイルストーンに追加
- `-w, --web`: 作成後にWebブラウザで開く

## Tips

- 複数行の本文には、改行や特殊文字を適切に処理するため、必ずheredoc構文を使用する
- 作成前にタイトルが空でないことを検証する
- ユーザーがラベル/担当者をカンマ区切りで提供した場合、複数のフラグに分割する
- 作成後、ユーザーが希望する場合はブラウザでIssueを開くことを提案する

## 例

```bash
gh issue create \
  --title "ユーザー認証機能の追加" \
  --body "$(cat <<'EOF'
以下の機能を含むユーザー認証を実装する必要があります：
- ログイン/ログアウト
- パスワードリセット
- セッション管理
EOF
)" \
  --label "enhancement" \
  --label "priority-high" \
  --assignee "username"
