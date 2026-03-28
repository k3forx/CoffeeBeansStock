# CLAUDE.md

このファイルはClaude Code (claude.ai/code)がこのリポジトリで作業する際のガイドラインを提供します。

## プロジェクト概要

Coffee Beans Stockは、複数ユーザー対応のコーヒー豆在庫管理モバイルアプリケーションです。Go（バックエンド）とReact Native + Expo（フロントエンド）で構築されます。

**主要な差別化要素**: 直近7日間の平均使用量に基づく消費ペース計算機能（「あと○日」表示）

## アーキテクチャ

### バックエンド (Go)

- **フレームワーク**: Chi ルーター（net/http互換、使用予定）
- **データベース**: PostgreSQL 15 + sqlc（型安全なクエリ生成）
- **マイグレーション**: golang-migrate（`backend/migrations/`）
- **アーキテクチャパターン**: クリーンアーキテクチャ + ドメイン駆動設計
- **認証**: JWT（アクセストークン: 1時間、リフレッシュトークン: 7日間、使用予定）

**重要**: `internal/database/` にはsqlcが生成したコードが含まれます。**手動編集厳禁**。修正する場合は `sql/schema.sql` または `sql/queries.sql` を編集し、`make sqlc-generate` で再生成してください。

### フロントエンド (React Native + Expo)

使用予定のライブラリ（計画中）:
- Expo Router（ファイルベースルーティング）
- Zustand（ローカル状態）+ TanStack Query v5（サーバー状態）
- NativeWind（React Native用Tailwind CSS）
- React Hook Form + Zodバリデーション

**現状**: 基本的なExpoセットアップのみ。

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

## API設計

ベースパス: `/api/v1`

**レスポンス形式**:
- 成功時: `{"success": true, "data": {...}}`
- エラー時: `{"success": false, "error": {"code": "ERROR_CODE", "message": "...", "details": {...}}}`

エンドポイント一覧はREADME.mdを参照。

## コーディング規約

### コミットメッセージ
```
<type>(<scope>): <subject>

<body>
```
タイプ: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`

## セキュリティ考慮事項

- 保護された全エンドポイントでJWT認証ミドルウェアを適用
- リソース所有権の検証: 常に`user_id`が認証ユーザーと一致するか確認
- パスワードはbcryptでハッシュ化（コスト: 10）
- SQLインジェクション対策: sqlcのパラメータ化クエリを使用
- シークレットは環境変数で管理（`.env`）

## 開発フェーズ

**MVP**（現在）: 認証、コーヒー豆CRUD、使用履歴、消費ペース計算
**フェーズ2**: 購入履歴、検索・フィルタリング、メール通知
**フェーズ3**: 画像アップロード、分析機能、テスト、CI/CD

## Issueワークフロー

Issue管理用のカスタムスキル:
- `/create-issue` - インタラクティブなissue作成
- `/implement-issue` - AI検証付き3フェーズ自動実装
