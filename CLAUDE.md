# CLAUDE.md

このファイルはClaude Code (claude.ai/code)がこのリポジトリで作業する際のガイドラインを提供します。

## プロジェクト概要

Coffee Beans Stockは、複数ユーザー対応のコーヒー豆在庫管理モバイルアプリケーションです。Go（バックエンド）とReact Native + Expo（フロントエンド）で構築され、ユーザーがコーヒー豆の追跡、使用履歴の記録、購入管理を行えます。

**主要な差別化要素**: 直近7日間の平均使用量に基づく消費ペース計算機能（「あと○日」表示）

## アーキテクチャ

### バックエンド (Go)

- **Go**: 1.26.0
- **モジュール**: `github.com/k3forx/CoffeeBeansStock/backend`
- **フレームワーク**: Chi ルーター（net/http互換、使用予定）
- **データベース**: PostgreSQL 15 + sqlc（型安全なクエリ生成）
- **マイグレーション**: golang-migrate（`backend/migrations/`）
- **アーキテクチャパターン**: クリーンアーキテクチャ + ドメイン駆動設計
  - `cmd/server/` - アプリケーションエントリーポイント
  - `internal/domain/` - ドメインエンティティとリポジトリインターフェース
  - `internal/application/` - アプリケーションサービス（ビジネスロジック）
  - `internal/services/` - ドメインサービス
  - `internal/database/` - sqlcが生成するデータベースコード（**手動編集厳禁**）
  - `internal/api/` - HTTPハンドラーとミドルウェア
  - `internal/auth/` - JWT認証（使用予定）
  - `internal/config/` - 設定管理
- **認証**: JWT（アクセストークン: 1時間、リフレッシュトークン: 7日間、使用予定）
- **ログ**: 標準ライブラリのslog

**現状**: ディレクトリ構造・SQLスキーマ・sqlc設定・Makefile・マイグレーションファイルは作成済み。Goソースコード（.goファイル）は未実装。

**重要**: `internal/database/` ディレクトリにはsqlcが生成したコードが含まれます。修正する場合は `sql/schema.sql` または `sql/queries.sql` を編集し、`sqlc generate` で再生成してください。

### フロントエンド (React Native + Expo)

- **プラットフォーム**: Expo SDK 54+（Managed Workflow）
- **言語**: TypeScript 5.9+
- **ナビゲーション**: Expo Router（ファイルベースルーティング、計画中）
- **状態管理**:
  - Zustand（ローカル状態、計画中）
  - TanStack Query v5（サーバー状態、計画中）
- **UI**: NativeWind（React Native用Tailwind CSS、計画中）
- **フォーム**: React Hook Form + Zodバリデーション（計画中）

**現状**: 基本的なExpoセットアップのみ。依存関係はREADME.md記載の手順でインストールが必要です。

### データベーススキーマ

- **users**: ユーザーアカウントと通知設定
- **coffee_beans**: コーヒー豆マスターデータ（名前、産地、焙煎度、現在在庫）
- **purchase_history**: 購入履歴（フェーズ2で実装）
- **usage_history**: 日々の使用記録、`usage_type`（'manual' または 'quick_button'）

全テーブルでUUID主キーを使用し、`deleted_at`によるソフトデリートに対応。

## 開発コマンド

### バックエンド

Makefileにより主要な操作はすべて `make` コマンドで実行可能（`backend/` ディレクトリで実行）:

```bash
cd backend

make run              # サーバー起動
make sqlc-generate    # sqlcコード生成
make migrate-up       # マイグレーション実行
make migrate-down     # マイグレーション1ステップ戻し
make test             # テスト実行（go test -v ./...）
make build            # バイナリビルド（bin/server）
make clean            # ビルド成果物削除
make help             # 利用可能なターゲット一覧
```

直接実行する場合:
```bash
# コードフォーマット
gofmt -w .

# テストカバレッジ
go test ./... -cover
```

**sqlcワークフロー**:
1. `sql/schema.sql` または `sql/queries.sql` を編集
2. `make sqlc-generate`（または `sqlc generate`）を実行
3. 生成されたコードが `internal/database/` に出力される
4. サービス層でインポートして使用

**注意**: sqlcのインストールが必要です: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

### フロントエンド

```bash
cd frontend

# 依存関係のインストール（未実施の場合 - package.json参照）
npm install

# 開発サーバー起動
npx expo start
# または: npm start

# プラットフォーム別起動
npx expo start --ios       # iOSシミュレータ
npx expo start --android   # Androidエミュレータ
npx expo start --web       # Web（開発時のみ）

# キャッシュクリア（トラブルシューティング）
npx expo start --clear
npx expo start --reset-cache

# 型チェック
npx tsc --noEmit
```

### データベース

```bash
# PostgreSQL起動
docker-compose up -d

# データベース状態確認
docker-compose ps

# マイグレーション実行（backend/ ディレクトリで）
make migrate-up

# マイグレーション1ステップ戻し
make migrate-down
```

**マイグレーションワークフロー**:
1. `backend/migrations/` に新しいマイグレーションファイルを作成（`NNNNNN_description.up.sql` / `.down.sql`）
2. `make migrate-up` で適用
3. スキーマ変更後は `sql/schema.sql` も同期し、`sqlc generate` を再実行

## API設計

ベースパス: `/api/v1`

### 認証エンドポイント
- `POST /auth/signup` - ユーザー登録
- `POST /auth/login` - ログイン（アクセストークン + リフレッシュトークンを返却）
- `POST /auth/refresh` - アクセストークン更新
- `GET /auth/me` - 現在のユーザー情報取得

### コーヒー豆エンドポイント
- `GET /coffee-beans` - ユーザーの豆一覧取得
- `POST /coffee-beans` - 新規豆登録
- `GET /coffee-beans/:id` - 豆詳細取得
- `PUT /coffee-beans/:id` - 豆更新
- `DELETE /coffee-beans/:id` - 豆削除（ソフトデリート）

### 使用履歴エンドポイント
- `GET /usage-history` - 使用履歴一覧取得
- `POST /usage-history` - 使用記録（手動入力）
- `POST /usage-history/quick` - 簡易記録（プリセット量）
- `DELETE /usage-history/:id` - 履歴削除

**レスポンス形式**:
- 成功時: `{"success": true, "data": {...}}`
- エラー時: `{"success": false, "error": {"code": "ERROR_CODE", "message": "...", "details": {...}}}`

## コーディング規約

### Go
- `gofmt` でフォーマット
- エクスポート名: パスカルケース（例: `CreateUser`）
- プライベート名: キャメルケース（例: `hashPassword`）
- 定数: 大文字スネークケース（例: `MAX_RETRIES`）
- SQL注釈: `sql/queries.sql` でsqlcコメントを使用（例: `-- name: GetUserByID :one`）

### TypeScript
- ESLint + Prettier を使用（設定は未定）
- 変数・関数: キャメルケース（例: `userName`）
- コンポーネント: パスカルケース（例: `BeanListScreen`）
- 型・インターフェース: パスカルケース（例: `CoffeeBean`、`AuthStore`）

### コミットメッセージ
```
<type>(<scope>): <subject>

<body>
```
タイプ: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`

例:
```
feat(auth): JWT認証機能を実装

- アクセストークンとリフレッシュトークンの発行
- 認証ミドルウェアの実装
- bcryptによるパスワードハッシュ化
```

## セキュリティ考慮事項

- 保護された全エンドポイントでJWT認証ミドルウェアを適用
- リソース所有権の検証: 常に`user_id`が認証ユーザーと一致するか確認
- パスワードはbcryptでハッシュ化（コスト: 10）
- SQLインジェクション対策: sqlcのパラメータ化クエリを使用
- 本番環境用のCORS設定が必要
- シークレットは環境変数で管理: `DB_PASSWORD`, `JWT_SECRET`, `SMTP_PASSWORD`（フェーズ2）

## 現在のプロジェクト状態

| コンポーネント | 状態 |
|---|---|
| DBスキーマ・マイグレーション | 作成済み（4テーブル） |
| sqlc設定・クエリ定義 | 作成済み |
| Docker Compose (PostgreSQL) | 作成済み |
| Makefile | 作成済み |
| バックエンド Goコード | 未実装（ディレクトリ構造のみ） |
| フロントエンド | Expoテンプレートのみ（カスタムコードなし） |
| CI/CD | 未設定 |

## 開発フェーズ

**MVP**（現在）: 認証、コーヒー豆CRUD、使用履歴、消費ペース計算
**フェーズ2**: 購入履歴、検索・フィルタリング、メール通知
**フェーズ3**: 画像アップロード、分析機能、テスト、CI/CD

## Issueワークフロー

このプロジェクトでは構造化されたGitHub issueテンプレート（`.github/ISSUE_TEMPLATE/`）を使用:
- **feature.yml** - 機能追加（フェーズ、実装領域、優先度、受け入れ条件）
- **bug_report.yml** - バグ報告
- **local_environment.yml** - ローカル環境セットアップ

Issue管理用のカスタムスキル:
- `/create-issue` - インタラクティブなissue作成
- `/implement-issue` - AI検証付き3フェーズ自動実装

## テスト戦略

- バックエンド: Goの標準testingパッケージ（`*_test.go`ファイル）
- フロントエンド: Jest + React Native Testing Library（設定保留中）
- APIエンドポイントの統合テスト
- サービス層のビジネスロジックのユニットテスト

## トラブルシューティング

**バックエンド**:
- データベース接続エラー: `.env`ファイルを確認し、PostgreSQLが起動しているか確認
- sqlcエラー: `sql/*.sql`のSQL構文を確認して再生成

**フロントエンド**:
- キャッシュ問題: `npx expo start --clear`
- Metro bundler問題: `npx expo start --reset-cache`
- 依存関係の問題: `rm -rf node_modules && npm install`
- ネイティブビルドが必要な場合: `npx expo prebuild`

## 参考リンク

- [sqlc Documentation](https://docs.sqlc.dev/)
- [Chi Router](https://github.com/go-chi/chi)
- [Expo Documentation](https://docs.expo.dev/)
- [TanStack Query](https://tanstack.com/query/latest/docs/react/overview)
- 詳細なプロジェクト要件とユーザーペルソナについてはREADME.mdを参照
