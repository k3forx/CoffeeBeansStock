---
name: implement-issue
shortDescription: GitHub Issueの内容を読み込み、実装計画を立案して自動実装し、受け入れ条件をチェックする
type: user-invocable
---

# GitHub Issue Implementation Skill

あなたはGitHub issueを自動実装するアシスタントです。このスキルは3段階のワークフローで動作します：

1. **Issue解析フェーズ（通常モード）**: Issue情報を取得・解析し、実装範囲を特定
2. **詳細計画フェーズ（Plan Mode）**: コードベースを深く探索し、詳細な実装計画を作成
3. **実装フェーズ（通常モード）**: コード実装、テスト、受け入れ条件チェック、コミット・PR作成

## Usage

```
/implement-issue <issue-number>
```

例:
```
/implement-issue 23
```

## Workflow Detection

**重要**: どのフェーズにいるかを自動検知して適切に動作します。

### フェーズ検知方法

システムリマインダーをチェックして現在のモードを判定:
- **Plan mode検知**: `<system-reminder>`に"Plan mode is active"が含まれている
- **通常モード**: Plan modeの指示がない

### フェーズ別の動作

```
Phase 1 (通常モード) → ユーザーに/planを提案 → Phase 2 (Plan mode) → ユーザー承認 → Phase 3 (通常モード)
```

## Phase 1: Issue Analysis (Regular Mode)

このフェーズは**通常モード**で実行されます。

### Steps

1. **Issue情報を取得**
   ```bash
   gh issue view <number> --json title,body,labels,number
   ```

2. **Issue内容を解析**
   - タイトルと本文を読み取る
   - 実装領域（フロントエンド/バックエンド/両方）を特定
   - 優先度を確認
   - 受け入れ条件を抽出（後のフェーズで使用）

3. **初期評価を提示**
   ```
   ## Issue #<number>: <title>

   **実装領域**: <frontend/backend/both>
   **優先度**: <High/Medium/Low>

   ### 概要
   <issue本文から抽出した概要>

   ### 受け入れ条件
   - [ ] 条件1
   - [ ] 条件2
   - [ ] 条件3
   ```

4. **Plan Modeへの移行を提案**
   ```
   詳細な実装計画を立案するため、Plan Modeに入ることをお勧めします。

   次のコマンドを実行してください:
   `/plan`

   その後、再度以下を実行してください:
   `/implement-issue <number>`
   ```

### ⚠️ このフェーズでやってはいけないこと

- コードベースの深い探索（Plan modeで行う）
- 実装の詳細な計画作成（Plan modeで行う）
- 実際のコード実装（Phase 3で行う）

## Phase 2: Detailed Planning (Plan Mode)

このフェーズは**Plan Mode**で実行されます。ユーザーが`/plan`コマンドを実行し、`/implement-issue <number>`を再実行したときに開始されます。

### Steps

1. **Plan Mode検知**
   - システムリマインダーで"Plan mode is active"を確認
   - Plan modeでない場合は、Phase 1に戻る

2. **Issue情報を再取得**
   ```bash
   gh issue view <number> --json title,body,labels,number
   ```

3. **コードベースの深い探索**

   実装領域に応じて以下を調査:

   **フロントエンド実装の場合:**
   - 既存のスクリーン構造を確認
   - 類似コンポーネントを検索
   - API呼び出しパターンを調査
   - ナビゲーション構造を理解

   ```bash
   # 例: スクリーンファイルの検索
   Glob: pattern="**/screens/**/*.tsx"

   # 例: API呼び出しパターンの検索
   Grep: pattern="useMutation|useQuery" path="frontend"
   ```

   **バックエンド実装の場合:**
   - 既存のハンドラー構造を確認
   - データベーススキーマを理解
   - sqlcクエリパターンを調査
   - ルーティング構造を確認

   ```bash
   # 例: ハンドラーの検索
   Glob: pattern="**/handlers/**/*.go"

   # 例: sqlcクエリの確認
   Read: file_path="backend/internal/db/query.sql"
   ```

4. **実装計画を作成**

   以下の形式で詳細な計画を作成:

   ```markdown
   # Implementation Plan for Issue #<number>

   ## Summary
   <実装の概要>

   ## Architecture Analysis
   <既存コードベースの分析結果>

   ## Implementation Steps

   ### 1. <ステップ名>
   **Files to create/modify:**
   - `path/to/file1.tsx` - 新規作成 - <説明>
   - `path/to/file2.go` - 変更 - <説明>

   **Changes:**
   - <具体的な変更内容>

   ### 2. <次のステップ>
   ...

   ## Testing Strategy
   - TypeScript: `cd frontend && npx tsc --noEmit`
   - Backend: `cd backend && go test ./...`

   ## Acceptance Criteria Mapping
   - [ ] 受け入れ条件1 → 実装ステップ1, 2で対応
   - [ ] 受け入れ条件2 → 実装ステップ3で対応
   - [ ] 受け入れ条件3 → 実装ステップ4で対応

   ## Risks and Considerations
   <リスクや注意点>
   ```

5. **ユーザー承認を待つ**

   計画を提示した後:
   ```
   上記の実装計画をご確認ください。

   承認いただける場合は、Plan Modeを終了してください。
   その後、実装フェーズが自動的に開始されます。
   ```

### ⚠️ このフェーズでやってはいけないこと

- 実際のコード実装（Phase 3で行う）
- ファイルの作成や編集（Phase 3で行う）
- git操作（Phase 3で行う）

## Phase 3: Implementation (Regular Mode)

このフェーズは**通常モード**（Plan mode終了後）で実行されます。

### Steps

1. **ブランチ作成**

   ```bash
   # Issue番号とタイトルからブランチ名を生成
   # 例: feature/issue-23-add-coffee-registration

   git checkout -b feature/issue-<number>-<slugified-title>
   ```

2. **コード実装**

   Phase 2で作成した計画に従って実装:
   - ファイルを作成/編集
   - 計画の各ステップを順次実行
   - 実装中に問題が発生したら調査・修正

3. **テスト実行**

   実装領域に応じてテストを実行:

   **フロントエンド:**
   ```bash
   cd frontend && npx tsc --noEmit
   ```

   **バックエンド:**
   ```bash
   cd backend && go test ./...
   ```

4. **受け入れ条件チェック（重要）**

   実装完了後、AIが自動的に受け入れ条件をチェックします。

   **チェック手順:**

   a. **Issue本文から受け入れ条件を抽出**
      - "## 受け入れ条件" セクションを探す
      - チェックボックスリスト `- [ ] ...` を抽出

   b. **実装内容の分析**
      - 変更されたファイルを Read ツールで確認
      - 追加/変更されたコードを解析

   c. **各条件をAIが判定**

      各受け入れ条件について:
      1. キーワードを抽出（例: "登録", "バリデーション", "遷移"）
      2. 実装ファイル内で関連コードを検索
      3. 機能の実装度を評価
      4. ステータスを決定:
         - ✅ **満たされている**: 機能が完全に実装されている
         - ⚠️ **部分的**: 機能は実装されているが不完全
         - ❌ **満たされていない**: 機能が実装されていない

   d. **結果を報告**

      ```markdown
      ## 受け入れ条件チェック結果

      ✅ ユーザーがコーヒー豆を登録できる
         → 実装確認: RegisterCoffeeScreen.tsx:45 にフォーム実装
         → API呼び出し: coffeeApi.ts:23 で POST /api/coffee

      ✅ 登録時にバリデーションが動作する
         → 実装確認: RegisterCoffeeScreen.tsx:67 でフォームバリデーション

      ⚠️ エラーメッセージが適切に表示される
         → 部分的に実装: API エラーは表示されるが、バリデーションエラーの UI が不足

      ❌ 登録完了後にリスト画面に遷移する
         → 未実装: navigation.goBack() のみ、リスト画面への直接遷移なし

      ---

      **判定サマリー:**
      - 満たされている: 2/4
      - 部分的: 1/4
      - 未実装: 1/4
      ```

   e. **追加実装の提案**

      満たされていない条件がある場合:
      ```
      未達の受け入れ条件があります。

      以下の対応が可能です:
      1. 追加実装を行う（未達の条件を実装します）
      2. このまま進める（現状のままコミット・PR作成）
      3. 中断する（変更を保留）

      どうしますか？
      ```

      ユーザーが「1. 追加実装」を選んだ場合:
      - 未達の条件を実装
      - 再度テスト実行
      - 受け入れ条件を再チェック

5. **変更をコミット**

   ```bash
   git add <modified-files>

   git commit -m "$(cat <<'EOF'
   feat: <title from issue>

   <summary of implementation>

   - Implemented feature X
   - Added validation for Y
   - Updated Z component

   Closes #<issue-number>

   Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
   EOF
   )"
   ```

6. **PRオプション**

   ユーザーに確認:
   ```
   実装が完了しました。

   次のアクションを選択してください:
   1. PRを作成する
   2. ブランチをpushのみ（PR作成は手動）
   3. ローカルのみ（pushしない）
   ```

   PRを作成する場合:
   ```bash
   git push -u origin feature/issue-<number>-<title>

   gh pr create \
     --title "feat: <title>" \
     --body "$(cat <<'EOF'
   ## Summary
   <implementation summary>

   ## Changes
   - Change 1
   - Change 2

   ## Related Issue
   Closes #<number>

   ## Acceptance Criteria
   <paste acceptance criteria check results>

   🤖 Generated with Claude Code
   EOF
   )"
   ```

## Error Handling

### Issue Not Found
```
Error: Issue #<number> が見つかりません。

Issue番号を確認してください。
既存のissueを確認: `gh issue list`
```

### GitHub CLI Not Authenticated
```
Error: GitHub CLI が認証されていません。

次のコマンドで認証してください:
`gh auth login`
```

### TypeScript Compilation Errors
```
TypeScript コンパイルエラーが発生しました:

<error details>

エラーを修正してから再度テストします...
```

### Test Failures
```
テストが失敗しました:

<test failure details>

このままコミットしますか？
1. はい（テスト失敗を承知でコミット）
2. いいえ（問題を修正する）
3. 中断（変更を保留）
```

### Git Operation Failures

**ブランチが既に存在:**
```
Warning: ブランチ feature/issue-<number> は既に存在します。

1. 既存ブランチを使用（上書きリスクあり）
2. 新しいブランチ名を使用（例: feature/issue-<number>-v2）
3. 中断

どうしますか？
```

**Push失敗:**
```
Error: git push に失敗しました。

考えられる原因:
- リモートブランチが更新されている
- 権限がない
- ネットワークエラー

再試行しますか？
```

### 受け入れ条件が存在しない

Issue本文に受け入れ条件セクションがない場合:
```
Warning: このissueには受け入れ条件が定義されていません。

受け入れ条件の自動チェックをスキップします。
実装内容はご自身で確認してください。
```

## Tips

### Best Practices

1. **Phase 1 → 2 → 3 の順序を守る**
   - 各フェーズには目的がある
   - 順序を飛ばすと品質が下がる

2. **Plan Modeでの探索を十分に行う**
   - 既存コードパターンを理解する
   - 一貫性のある実装につながる

3. **受け入れ条件を基準にする**
   - 実装が受け入れ条件を満たすことを最優先
   - チェック結果を真剣に受け止める

4. **テスト失敗を無視しない**
   - テストが失敗したら原因を調査
   - 品質を担保する

5. **コミットメッセージを丁寧に**
   - 何を実装したか明確に
   - Issue番号を必ず含める

### Common Gotchas

1. **Plan Modeの検知ミス**
   - システムリマインダーを必ずチェック
   - 誤ったフェーズで動作すると混乱の原因

2. **受け入れ条件の抽出失敗**
   - Issue本文のフォーマットが異なる場合がある
   - "受け入れ条件"だけでなく"完了条件"なども検索

3. **大規模な実装**
   - 複数機能にまたがる場合、サブタスクへの分割を提案
   - 一度に全てを実装しようとしない

4. **依存関係**
   - 他のissueに依存している場合、先に完了させる必要がある
   - Issue本文の"Related Issues"を確認

### Recommendations

1. **Issue番号を間違えない**
   - `gh issue list` で確認してから実行

2. **Plan Modeを活用する**
   - 複雑な実装ほどPlan Modeの価値が高い
   - 計画に時間をかけることで実装がスムーズになる

3. **受け入れ条件を明確に**
   - Issueを作成する段階で受け入れ条件を詳細に書く
   - 自動チェックの精度が上がる

4. **段階的にコミット**
   - 大きな変更は複数コミットに分ける
   - 履歴が追いやすくなる

## GitHub CLI Commands Reference

### Issue Operations
```bash
# Issue情報取得（JSON形式）
gh issue view <number> --json title,body,labels,assignees,milestone,number

# Issue情報取得（人間が読みやすい形式）
gh issue view <number>

# Issueにコメント追加
gh issue comment <number> --body "実装完了しました"

# Issueラベル編集
gh issue edit <number> --add-label "status: implemented"
```

### Git Operations
```bash
# ブランチ作成 & 切り替え
git checkout -b feature/issue-<number>-<description>

# 現在のブランチ確認
git branch --show-current

# 変更ファイル確認
git status

# ファイルをステージング
git add <files>

# コミット作成
git commit -m "message"

# リモートへpush
git push -u origin <branch-name>
```

### PR Operations
```bash
# PR作成
gh pr create --title "title" --body "body"

# PR一覧
gh pr list

# PR詳細
gh pr view <number>
```

## Example Workflows

### Example 1: シンプルな機能追加（フロントエンド）

Issue #23: ボタンコンポーネントの追加

**Phase 1 (通常モード):**
```
User: /implement-issue 23
AI: [Issue情報を取得・表示]
AI: "Plan Modeで詳細計画を立案しますか？ → /plan を実行してください"
```

**Phase 2 (Plan Mode):**
```
User: /plan
User: /implement-issue 23
AI: [コードベース探索]
AI: [実装計画を作成・提示]
AI: "計画をご確認ください。承認いただければ、Plan Modeを終了してください"
User: [承認]
```

**Phase 3 (通常モード):**
```
AI: [ブランチ作成]
AI: [ファイル作成・実装]
AI: [テスト実行]
AI: [受け入れ条件チェック]
AI: "✅ すべての受け入れ条件が満たされています"
AI: [コミット]
AI: "PRを作成しますか？"
User: はい
AI: [PR作成]
AI: "実装完了！PR URL: https://github.com/..."
```

### Example 2: 複雑な機能（フロントエンド + バックエンド）

Issue #45: コーヒー豆在庫管理機能

**Phase 1:**
```
AI: "実装領域: 両方（フロントエンド + バックエンド）"
AI: "複雑な実装のため、Plan Modeでの詳細計画を強く推奨します"
```

**Phase 2:**
```
AI: [フロントエンド既存パターン調査]
AI: [バックエンドAPI構造調査]
AI: [データベーススキーマ確認]
AI: [詳細な実装計画作成（10ステップ）]
```

**Phase 3:**
```
AI: [フロントエンド実装]
AI: [バックエンドAPI実装]
AI: [テスト実行]
AI: [受け入れ条件チェック]
AI: "⚠️ 1つの条件が部分的にのみ実装されています"
AI: "追加実装しますか？"
User: はい
AI: [追加実装]
AI: [再チェック]
AI: "✅ すべての条件を満たしました"
AI: [コミット & PR]
```

### Example 3: 受け入れ条件未達のケース

Issue #67: ログイン機能

**Phase 3での受け入れ条件チェック:**
```
AI: "受け入れ条件チェック結果:
     ✅ ログインフォームが表示される
     ✅ 認証APIが呼ばれる
     ⚠️ エラーメッセージが表示される（部分的）
     ❌ ログイン成功後にダッシュボードに遷移する（未実装）"

AI: "未達の条件があります。どうしますか？
     1. 追加実装を行う
     2. このまま進める
     3. 中断する"

User: 1

AI: [未達の条件を実装]
AI: [再テスト & 再チェック]
AI: "✅ すべての条件を満たしました"
```

## Advanced Topics

### カスタム受け入れ条件フォーマット

標準的なフォーマット以外の場合:
```markdown
# Issue本文の例（非標準）

完了条件：
* 機能Aが動作する
* 機能Bが動作する
```

AIは以下も検出:
- "完了条件"
- "Definition of Done"
- "DoD"
- 箇条書きリスト（`-` または `*`）

### 並行開発の考慮

複数人が同じissueに取り組む可能性:
```
Warning: このissueには既に以下のPRがあります:
- PR #123 by @user1 (Open)

続行しますか？ (競合の可能性があります)
```

### Issue依存関係の自動検出

Issue本文に "Blocks: #123" や "Depends on: #456" がある場合:
```
Warning: このissueは #123 に依存しています。

#123 の状態: Open

先に #123 を完了させることを推奨します。
続行しますか？
```

## Troubleshooting

### "Plan mode is active" が検出されない

原因: システムリマインダーの確認ミス

解決策:
- システムリマインダーを再確認
- ユーザーに「/plan コマンドを実行しましたか？」と確認

### 受け入れ条件の抽出に失敗

原因: Issue本文のフォーマットが想定外

解決策:
- 複数のキーワードで検索（"受け入れ条件", "完了条件", "Acceptance Criteria"）
- チェックボックス以外のリスト形式も試す
- 抽出できない場合はユーザーに確認

### テストが常に失敗する

原因: 環境の問題、依存関係の不足

解決策:
```bash
# 依存関係を再インストール
cd frontend && npm install
cd backend && go mod tidy

# 再テスト
```

### ブランチ作成に失敗

原因: 未コミットの変更がある

解決策:
```bash
# 現在の変更を確認
git status

# 変更をstash
git stash

# ブランチ作成
git checkout -b feature/issue-<number>

# 必要に応じてstash適用
git stash pop
```

## Summary

このスキルは以下を実現します:

1. **3段階ワークフロー**: 解析 → 計画 → 実装
2. **Plan Mode統合**: 詳細な計画立案で品質向上
3. **受け入れ条件の自動チェック**: AI判定で実装完了度を可視化
4. **エンドツーエンドの自動化**: Issue → 実装 → PR まで一貫したフロー

使用することで:
- ✅ 開発速度が向上
- ✅ 実装品質が担保
- ✅ 受け入れ条件の見落としを防止
- ✅ 一貫性のあるコード

それでは `/implement-issue <number>` を実行してください！
