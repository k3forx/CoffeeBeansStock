# CLAUDE.md

## プロジェクト概要

Coffee Beans Stockは、複数ユーザー対応のコーヒー豆在庫管理モバイルアプリケーションです。Go（バックエンド）とReact Native + Expo（フロントエンド）で構築されます。

**主要な差別化要素**: 直近7日間の平均使用量に基づく消費ペース計算機能（「あと○日」表示）

## 重要な注意事項

`internal/database/` にはsqlcが生成したコードが含まれます。**手動編集厳禁**。修正する場合は `sql/schema.sql` または `sql/queries.sql` を編集し、`make sqlc-generate` で再生成してください。

## 開発コマンド

### Backend (`backend/`)

| コマンド | 用途 |
|---------|------|
| `make run` | サーバー起動 |
| `make test` | テスト実行（`-v -race`） |
| `make lint` | golangci-lint実行 |
| `make build` | バイナリビルド |
| `make sqlc-generate` | sqlcコード再生成 |
| `make migrate-up` | マイグレーション適用 |
| `make migrate-down` | マイグレーション1ステップ戻す |
| `make generate-api` | OpenAPIからGoの型を生成 |

### Frontend (`frontend/`)

| コマンド | 用途 |
|---------|------|
| `npm test` | Jestテスト実行 |
| `npm run lint` | ESLint実行 |
| `npm run typecheck` | TypeScript型チェック |
| `npm run generate:types` | OpenAPIからTypeScript型を生成 |

## マイグレーションワークフロー

1. `backend/migrations/` にファイルを作成（`NNNNNN_description.up.sql` / `.down.sql`）
2. `make migrate-up` で適用
3. スキーマ変更後は `sql/schema.sql` も同期し、`make sqlc-generate` を再実行

## OpenAPI型生成ワークフロー

APIスキーマ（`backend/api/openapi.yaml`）を変更した場合、両方の型生成を実行すること:

1. `cd backend && make generate-api` — Goの型を再生成
2. `cd frontend && npm run generate:types` — TypeScriptの型を再生成

## Pre-commitフック

lefthookにより、コミット時に以下が自動実行される:

- **backend**: `golangci-lint run`（変更ファイル対象）
- **frontend**: `npx eslint`（変更ファイル対象）+ `npx tsc --noEmit`

コミット前にこれらが通ることを確認すること。

## ブランチ命名規則

```
feature/issue-<number>-<description>
```

例: `feature/issue-67-consumption-pace-prediction`

## コーディング規約

### コミットメッセージ

```
<type>(<scope>): <subject>

<body>
```

タイプ: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`

### テスト方針

**Backend**:
- `go-cmp` でアサーション（testify禁止）
- テーブル駆動テスト（`map[string]struct{}`形式）
- テストケース名は日本語
- 詳細は `.claude/rules/backend-testing.md` 参照

**Frontend**:
- Jest + `@testing-library/react-native`
- テスト名は日本語
- スナップショットテスト不使用
- 詳細は `.claude/rules/frontend-testing.md` 参照

### アーキテクチャ（バックエンド）

軽量DDD（ThreeDotsLabs「DDD Lite」）を採用。domain層は外部依存なし（標準ライブラリ+UUIDのみ）。詳細は `.claude/rules/backend-ddd.md` 参照。
