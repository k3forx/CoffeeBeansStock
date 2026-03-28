---
paths:
  - "backend/**/*.go"
---

# バックエンド 軽量DDD アーキテクチャ方針

ThreeDotsLabs「DDD Lite」に基づく戦術的パターンのみ採用。Go idiom重視。

## レイヤー構成と責務

```
domain/          ビジネスルールの中核。外部依存なし（標準ライブラリ+UUIDのみ）
  user/          User集約（エンティティ + Repository interface）
  coffeebean/    CoffeeBean集約（エンティティ + VO + Repository interface）
  auth/          TokenManager interface
  unitofwork/    UnitOfWork + Store interface
  errors.go      共通ドメインエラー（ErrNotFound, ErrForbidden等）+ ValidationError
  quantity.go    共通VO（複数集約で利用）

services/        アプリケーション層。ユースケースの組み立て。ドメインロジックは書かない
repository/      domain層インターフェースの実装（SQLC経由）
api/handlers/    HTTPリクエスト/レスポンス変換。domain型↔JSON変換はここで行う
api/middleware/  認証ミドルウェア等
auth/            JWTManager実装（domain/auth.TokenManager を満たす）
database/        SQLC生成コード（手動編集禁止）
```

## エンティティ設計パターン

- **unexportedフィールド** + publicアクセサ（getter）で不変条件を保護
- **ファクトリ関数** `New*()`：生成時にバリデーション実行、`error`を返す
- **復元関数** `Reconstruct()`：DB→ドメインモデル変換用。バリデーションなし
- ビジネスルール違反は `domain.ValidationError` / `domain.ValidationErrors` を返す
- IDはドメイン層で `uuid.New()` 生成し、DBに渡す

## Value Object 設計パターン

- 不変（フィールドはunexported、変更メソッドは新インスタンスを返す）
- ファクトリ `New*()` で不変条件を保証（例: `Stock: 0≤x≤50000`, `RoastLevel: 8段階enum`）
- DB復元用に `Reconstruct*()` を提供（バリデーションなし）
- 共通VOは `domain/` ルート、集約固有VOは集約パッケージ内に配置

## Repository interface

- **domain層で定義**（`user/repository.go`, `coffeebean/repository.go`）
- `repository/` パッケージで実装（`database.DBTX` を受け取り Pool/Tx 両対応）
- 返り値はドメイン型（`*user.User`, `*coffeebean.CoffeeBean`）
- `pgx.ErrNoRows` → `domain.ErrNotFound` に変換

## UnitOfWork パターン

- `unitofwork.UnitOfWork.RunInTx(ctx, func(store Store) error)` で集約間トランザクション管理
- `Store` 経由でtx内Repositoryを取得（構造的にtx漏れを防止）
- 単一集約操作ではUoW不要（Repository単体で十分）

## エラーハンドリングフロー

```
domain層: domain.ErrNotFound, domain.ErrForbidden, domain.ValidationError 等を返す
  ↓
service層: そのまま伝播（エラー変換しない）
  ↓
handler層: errors.Is/errors.As で判別し HTTPステータスに変換
  domain.ErrNotFound       → 404
  domain.ErrForbidden      → 403
  domain.ValidationError   → 400
  domain.ErrInvalidToken   → 401
  domain.ErrExpiredToken   → 401
```

## 導入しないもの（意図的判断）

- Email VO / PasswordHasher — 匿名認証のため現在不要
- Domain Events / CQRS — 小規模アプリに過剰
- BeanName等の単純なVO — 集約内ファクトリでのバリデーションで十分
- Bounded Context分割 — 単一コンテキストで運用
