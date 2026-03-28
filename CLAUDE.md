# CLAUDE.md

このファイルはClaude Code (claude.ai/code)がこのリポジトリで作業する際のガイドラインを提供します。

## プロジェクト概要

Coffee Beans Stockは、複数ユーザー対応のコーヒー豆在庫管理モバイルアプリケーションです。Go（バックエンド）とReact Native + Expo（フロントエンド）で構築されます。

**主要な差別化要素**: 直近7日間の平均使用量に基づく消費ペース計算機能（「あと○日」表示）

## 重要な注意事項

`internal/database/` にはsqlcが生成したコードが含まれます。**手動編集厳禁**。修正する場合は `sql/schema.sql` または `sql/queries.sql` を編集し、`make sqlc-generate` で再生成してください。

## 開発コマンド

バックエンドは `backend/` ディレクトリで `make help` を実行するとターゲット一覧を確認可能。

### sqlcワークフロー
1. `sql/schema.sql` または `sql/queries.sql` を編集
2. `make sqlc-generate` を実行
3. 生成されたコードが `internal/database/` に出力される

### マイグレーションワークフロー
1. `backend/migrations/` にファイルを作成（`NNNNNN_description.up.sql` / `.down.sql`）
2. `make migrate-up` で適用
3. スキーマ変更後は `sql/schema.sql` も同期し、`make sqlc-generate` を再実行

## コーディング規約

### コミットメッセージ
```
<type>(<scope>): <subject>

<body>
```
タイプ: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`

## Issueワークフロー

Issue管理用のカスタムスキル:
- `/create-issue` - インタラクティブなissue作成
- `/implement-issue` - AI検証付き3フェーズ自動実装
