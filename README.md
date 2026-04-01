# Coffee Beans Stock - プロジェクトガイド

## プロジェクト概要

複数ユーザーが利用できるコーヒー豆の在庫管理モバイルアプリケーション。ユーザーは自分が所有するコーヒー豆の情報（名前、産地、焙煎度、購入情報、在庫量）を管理し、使用履歴や購入履歴を記録できます。在庫が設定した閾値を下回った場合にはメールで通知を受け取ることができます。

### 主要機能
- ユーザー認証（JWT）
- コーヒー豆のCRUD操作
- 在庫管理
- 使用履歴記録
- 購入履歴記録（フェーズ2）
- 検索・フィルタリング（フェーズ2）
- メール通知（フェーズ2）

## ユーザー像とユースケース

### ペルソナ（ターゲットユーザー）

**名前**: 田中太郎（仮）
**属性**: 30代、コーヒー愛好家、在宅勤務
**コーヒー習慣**:
- 1日2-3杯のコーヒーを自宅で淹れる
- 常に3-5種類の豆を保有し、気分や残量で使い分け
- 焙煎後2-3週間以内に使い切りたい

**課題**:
- 豆の在庫が把握しづらく、「切らしてしまった」「古い豆が残っている」という事態が発生
- 購入したコーヒー豆の総コストが見えず、家計管理が難しい
- 複数の豆を管理する手間を最小限にしたい

**ニーズ**:
- 在庫が少なくなったら通知が欲しい
- 各豆の「あと何日で使い切れるか」を知りたい
- 簡単に記録できる（手間をかけたくない）

### 主要なユースケース

#### UC-01: コーヒーを淹れる
1. ユーザーがコーヒーを淹れる前にアプリを開く
2. 豆の一覧画面で各豆の状態を確認
   - 現在の在庫量（グラム）
   - 消費ペース（「あと5日」など）
   - 鮮度警告（購入から○○日経過）
3. 気分または残量を見て、今日使う豆を選択
4. コーヒーを淹れる
5. 使用後、アプリで使用量（グラム）を記録
6. アプリが在庫を自動更新し、必要に応じて警告を表示

#### UC-02: 在庫補充タイミングの把握
1. 豆を使用した際、在庫が閾値（デフォルト100g）を下回る
2. **アプリ内で警告を表示**（例: 「エチオピア イルガチェフェの在庫が少なくなっています」）
3. **メールで通知**（閾値を下回った時点で1回のみ送信）
4. ユーザーが通知を見て、次回の買い物時に補充

#### UC-03: 新しい豆を追加
1. コーヒー豆を購入
2. アプリで新規登録
   - 名前（例: エチオピア イルガチェフェ）
   - 産地、焙煎度
   - 購入日、購入価格、購入店舗
   - 初期在庫量（例: 200g）
3. 一覧画面に追加され、以降の使用履歴記録が可能になる

#### UC-04: 過去の履歴確認
- 使用履歴画面で「いつ・どの豆を・何グラム使ったか」を確認
- 各豆の詳細画面で累計消費量や購入履歴を確認
- **詳細な分析機能は不要**（シンプルに履歴が見られれば十分）

### ユーザージャーニー

```
[朝 - コーヒータイム]
1. アプリ起動
   ↓
2. 豆一覧を確認
   - エチオピア: 120g（あと4日）
   - ブラジル: 50g（⚠️ あと1.5日）← 警告表示
   - コロンビア: 180g（あと6日）
   ↓
3. 今日はエチオピアを選択（気分で）
   ↓
4. コーヒーを淹れる（15g使用）
   ↓
5. 使用記録画面で「エチオピア」を選び、「15g」を入力
   ↓
6. 在庫が自動更新（120g → 105g）

[翌朝]
1. アプリ起動
   ↓
2. ブラジルが閾値以下 → 「在庫が少なくなっています」と表示
   ↓
3. 同時にメール通知を受信
   ↓
4. 週末に買い物リストに追加
```

### UI/UX要件（重要）

#### 🎯 最優先事項
1. **豆一覧のわかりやすさ**
   - 在庫量を大きく表示
   - 「あと○日」の表示（消費ペース計算）
   - 警告アイコン（閾値以下の豆）

2. **豆選択の簡単さ**
   - **課題**: リストから豆を探すのが面倒
   - **対策案**:
     - 最近使った豆を上位表示
     - お気に入り機能（ピン留め）
     - 豆の画像やカラーラベル（視覚的に識別しやすく）
     - 検索・フィルタ機能（フェーズ2）

3. **記録の手軽さ**
   - 使用量入力は数値キーパッドで素早く入力
   - デフォルト値の設定（例: いつも15gなら初期値15g）
   - 日付は「今日」がデフォルト

#### 📊 表示優先度

**豆一覧画面で表示すべき情報（優先順）**:
1. 現在の在庫量（グラム） + 警告マーク
2. 消費ペース（あと何日）← **新規実装が必要**
3. 豆の名前
4. 購入日からの経過日数（鮮度）

**詳細画面で表示する情報**:
- 基本情報（名前、産地、焙煎度、メモ）
- 購入情報（購入日、価格、店舗）
- 累計消費量、1杯あたりのコスト
- 使用履歴（日付、使用量）

### 新規実装が必要な機能

#### ✨ 消費ペース計算機能（必須）
- **計算ロジック**:
  - 直近7日間の平均使用量を算出
  - 現在の在庫量 ÷ 1日平均使用量 = 残り日数
  - 例: 在庫120g、直近7日で平均20g/日 → あと6日
- **表示**:
  - 「あと6日」のような自然言語表示
  - データが不足する場合は「データ不足」と表示
- **実装場所**:
  - バックエンド: サービス層で計算ロジック実装
  - フロントエンド: 一覧画面・詳細画面に表示

#### 📧 通知機能の詳細化
- **アプリ内警告**: 使用記録時にリアルタイムで警告表示
- **メール通知**: 閾値以下になった時点で1回送信（重複防止）

## 市場調査と差別化戦略

### 競合状況
- **日本市場**: 記録・レビュー系アプリが主流（コーヒー手帳、コーヒーノート等）。在庫管理特化型はほぼ不在。
- **海外市場**: Beanconqueror（在庫追跡あり、オープンソース）、BeanVault、Bean Tracker等が存在。

### 差別化ポイント（MVPで実現）
1. **在庫管理特化**（日本市場で唯一）
2. **消費ペース計算**（競合にない独自機能 - 「あと○日」表示）
3. **簡易ボタン記録**（ワンタップで使用記録、競合は手動入力のみ）
4. **複数ユーザー対応**（JWT認証、家族・シェアハウス利用）
5. **日本語ネイティブ対応**（Beanconquerorは英語のみ）

### 競合比較（簡易版）
| 機能 | 本アプリ | Beanconqueror | 日本の既存アプリ |
|------|---------|---------------|---------------|
| 在庫管理 | ✅ 中心 | ✅ | ❌ ほぼなし |
| 消費ペース計算 | ✅ 独自 | ❌ | ❌ |
| 簡易ボタン記録 | ✅ | ❌ 手動のみ | 一部 |
| メール通知 | フェーズ2 | ❌ | ❌ |
| 日本語UI | ✅ | ❌ | ✅ |

### 優先順位の整理

#### MVP（現行設計）
- ユーザー認証
- コーヒー豆CRUD
- 使用履歴記録
- **消費ペース計算機能**（追加）
- **豆一覧UI改善**（最近使った豆を上位表示、視覚的識別）
- **アプリ内警告表示**（追加）
- 基本的なUI

#### フェーズ2
- 購入履歴管理
- 検索・フィルタリング
- メール通知機能
- ユーザー設定画面
- お気に入り機能

#### フェーズ3
- 画像・カラーラベル機能
- 詳細な分析機能（不要との声あり → 優先度低）
- テストカバレッジ向上
- UI/UX改善
- パフォーマンス最適化
- CI/CD整備

## 技術スタック

### バックエンド
- **言語**: Go 1.26.0
- **フレームワーク**: Chi (標準ライブラリベースの軽量ルーター)
- **データベース**: PostgreSQL 15 + sqlc (型安全なSQLコード生成) + pgx/v5
- **認証**: JWT (golang-jwt/jwt/v5)
- **設定管理**: 標準ライブラリ (os.Getenv) + godotenv
- **ログ**: slog (Go標準ライブラリ)
- **メール送信**: gomail（フェーズ2）

#### 技術選定の理由
- **Chi**: net/http完全互換、シンプル、豊富なミドルウェアエコシステム
- **sqlc**: SQL直接記述で型安全なGoコード生成、パフォーマンス良好
- **slog**: 標準ログライブラリ、外部依存なし、構造化ログ対応

### フロントエンド
- **プラットフォーム**: Expo SDK 50+（Managed Workflow）
- **言語**: TypeScript 5.3+
- **UI**:
  - スタイリング: NativeWind（Tailwind CSS for React Native）
  - 基本コンポーネント: React Native標準 + カスタムコンポーネント
- **状態管理**:
  - 認証・ローカル状態: Zustand + zustand/middleware/persist
  - サーバー状態: TanStack Query v5（React Query）
- **ナビゲーション**: Expo Router（React Navigation v6ベース）
- **フォーム**: React Hook Form + Zod
- **HTTP通信**: Axios（TanStack QueryのqueryFn内で使用）
- **ストレージ**: AsyncStorage（Zustandの永続化用）
- **テスト**: Jest + React Native Testing Library

#### 技術選定の理由
- **Expo**: MVP開発速度優先、必要な機能は全てExpoでカバー可能、OTA更新対応
- **NativeWind**: Tailwind CSS記法で高速スタイリング、Web版への展開も容易
- **TanStack Query**: サーバー状態管理の自動化（キャッシュ、リフェッチ、楽観的更新）
- **Zod**: TypeScriptと統合された型安全なバリデーション、バックエンド型定義との共有可能
- **Expo Router**: ファイルベースルーティング、型安全なナビゲーション

### インフラ
- **開発環境**: Docker Compose
- **データベース**: PostgreSQL 15

## プロジェクト構造

```
CoffeeBeansStock/
├── backend/                    # Go APIサーバー
│   ├── cmd/server/
│   │   └── main.go            # エントリーポイント
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/      # APIハンドラー
│   │   │   ├── middleware/    # 認証、CORS等
│   │   │   └── routes.go      # Chiルート定義
│   │   ├── database/          # sqlcが生成するコード
│   │   │   ├── db.go          # sqlc生成: DB接続インターフェース
│   │   │   ├── models.go      # sqlc生成: モデル定義
│   │   │   └── queries.sql.go # sqlc生成: クエリ実装
│   │   ├── services/          # ビジネスロジック
│   │   ├── auth/              # JWT実装
│   │   ├── email/             # メール送信（フェーズ2）
│   │   └── config/            # 設定管理
│   ├── sql/                   # sqlcのソースSQL
│   │   ├── schema.sql         # テーブル定義（DDL）
│   │   └── queries.sql        # クエリ定義
│   ├── migrations/            # DBマイグレーション（golang-migrate等）
│   ├── tests/
│   ├── sqlc.yaml              # sqlc設定ファイル
│   ├── .env.example           # 環境変数テンプレート
│   ├── go.mod
│   └── Makefile
├── frontend/                  # React Native (Expo)
│   ├── app/                  # Expo Router（ファイルベースルーティング）
│   │   ├── (auth)/           # 認証グループ
│   │   │   ├── login.tsx
│   │   │   └── signup.tsx
│   │   ├── (tabs)/           # タブナビゲーション
│   │   │   ├── index.tsx     # 豆一覧画面
│   │   │   ├── history.tsx   # 使用履歴画面
│   │   │   └── profile.tsx   # プロフィール画面
│   │   ├── beans/
│   │   │   ├── [id].tsx      # 豆詳細画面
│   │   │   └── new.tsx       # 豆新規登録画面
│   │   └── _layout.tsx       # ルートレイアウト
│   ├── src/
│   │   ├── api/              # APIクライアント（Axios + TanStack Query）
│   │   ├── components/       # 共通コンポーネント
│   │   ├── store/            # 状態管理（Zustand）
│   │   ├── types/            # TypeScript型定義
│   │   ├── utils/            # ユーティリティ関数
│   │   └── hooks/            # カスタムフック（TanStack Query等）
│   ├── assets/               # 画像・フォント等
│   ├── app.json              # Expo設定
│   ├── tailwind.config.js    # NativeWind設定
│   └── package.json
├── docker-compose.yml         # ローカル開発環境
├── CLAUDE.md                  # このファイル
└── README.md
```

## 認証設計

### JWT認証フロー

1. **Access Token**（短命）
   - 有効期限: 1時間
   - 用途: API認証
   - 保存場所: AsyncStorage

2. **Refresh Token**（長命）
   - 有効期限: 7日間
   - 用途: Access Token更新
   - 保存場所: AsyncStorage

### トークン構造

```go
type AccessTokenClaims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}
```

### 認証フロー

1. ユーザーがログイン → Access Token + Refresh Tokenを取得
2. APIリクエスト時に`Authorization: Bearer {access_token}`ヘッダーを付与
3. Access Token期限切れ（401エラー）→ 自動的にRefresh Tokenでトークン更新
4. Refresh Token期限切れ → ログイン画面へリダイレクト

## API設計

### エンドポイント一覧（MVP）

#### 認証
```
POST   /api/v1/auth/signup       # ユーザー登録
POST   /api/v1/auth/login        # ログイン
POST   /api/v1/auth/refresh      # トークンリフレッシュ
GET    /api/v1/auth/me           # 現在のユーザー情報
```

#### コーヒー豆管理
```
GET    /api/v1/coffee-beans      # 一覧取得
POST   /api/v1/coffee-beans      # 新規登録
GET    /api/v1/coffee-beans/:id  # 詳細取得
PUT    /api/v1/coffee-beans/:id  # 更新
DELETE /api/v1/coffee-beans/:id  # 削除（ソフトデリート）
```

#### 使用履歴
```
GET    /api/v1/usage-history           # 一覧取得
POST   /api/v1/usage-history           # 手動記録
POST   /api/v1/usage-history/quick     # 簡易ボタン記録
DELETE /api/v1/usage-history/:id       # 削除
```

### レスポンス形式

**成功時:**
```json
{
  "success": true,
  "data": { ... }
}
```

**エラー時:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": { ... }
  }
}
```

## データベース設計

### テーブル構成

#### users
- id (UUID, PK)
- email (VARCHAR, UNIQUE)
- password_hash (VARCHAR)
- name (VARCHAR)
- low_stock_threshold (INTEGER, デフォルト: 100g)
- created_at, updated_at, deleted_at

#### coffee_beans
- id (UUID, PK)
- user_id (UUID, FK → users.id)
- name (VARCHAR)
- origin (VARCHAR) - 産地
- roast_level (VARCHAR) - 焙煎度
- current_stock (INTEGER) - 現在の在庫（グラム）
- purchase_date (DATE)
- purchase_price (DECIMAL)
- purchase_store (VARCHAR)
- notes (TEXT)
- created_at, updated_at, deleted_at

#### usage_history
- id (UUID, PK)
- coffee_bean_id (UUID, FK → coffee_beans.id)
- user_id (UUID, FK → users.id)
- usage_date (DATE)
- quantity (INTEGER) - 使用量（グラム）
- usage_type (VARCHAR) - 'manual' or 'quick_button'
- notes (TEXT)
- created_at

## 開発規約

### コーディング規約

**Go:**
- `gofmt`でフォーマット
- 変数名: キャメルケース（例: `userName`）
- 定数名: 大文字スネークケース（例: `MAX_RETRIES`）
- エクスポート関数: 大文字始まり（例: `CreateUser`）
- プライベート関数: 小文字始まり（例: `hashPassword`）

**TypeScript:**
- ESLint + Prettierでフォーマット
- 変数名: キャメルケース（例: `userName`）
- コンポーネント名: パスカルケース（例: `BeanListScreen`）
- 型名: パスカルケース（例: `CoffeeBean`）
- インターフェース名: パスカルケース（例: `AuthStore`）

### コミットメッセージ規約

```
<type>(<scope>): <subject>

<body>
```

**Type:**
- `feat`: 新機能
- `fix`: バグ修正
- `refactor`: リファクタリング
- `docs`: ドキュメント
- `test`: テスト
- `chore`: その他（依存関係更新等）

**例:**
```
feat(auth): JWT認証機能を実装

- Access TokenとRefresh Tokenの発行
- 認証ミドルウェアの実装
- パスワードハッシュ化（bcrypt）
```

### ブランチ戦略

- `main`: 本番環境
- `develop`: 開発環境
- `feature/{feature-name}`: 機能開発
- `fix/{bug-name}`: バグ修正

## セキュリティ考慮事項

### 認証・認可
- 全保護エンドポイントで認証ミドルウェア適用
- リソースオーナーシップ確認（user_id照合）
- パスワードはbcryptでハッシュ化（cost: 10）

### API
- CORS設定（許可されたオリジンのみ）
- 入力バリデーション（すべてのエンドポイント）
- SQLインジェクション対策（プリペアドステートメント）

### 環境変数
- `.env`ファイルで管理（gitignore）
- 必須: DB_PASSWORD, JWT_SECRET, SMTP_PASSWORD（フェーズ2）

## 実装フェーズ

### MVP（Week 1-3）
- ユーザー認証
- コーヒー豆CRUD
- 使用履歴記録
- 基本的なUI

### フェーズ2（Week 4-5）
- 購入履歴管理
- 検索・フィルタリング
- メール通知機能
- ユーザー設定画面

### フェーズ3（Week 6）
- テストカバレッジ向上
- UI/UX改善
- パフォーマンス最適化
- CI/CD整備

## 開発環境セットアップ

### バックエンド

```bash
# PostgreSQL起動
docker-compose up -d

# 依存関係インストール
cd backend
go mod download

# sqlcインストール（初回のみ）
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# sqlcでコード生成
make sqlc-generate
# または: sqlc generate

# マイグレーション実行
make migrate-up

# サーバー起動
make run
```

**Note**: sqlcは`sql/schema.sql`と`sql/queries.sql`から`internal/database/`配下にGoコードを自動生成します。SQLファイルを変更した際は必ず`sqlc generate`を実行してください。

### フロントエンド

```bash
# プロジェクト作成（初回のみ）
npx create-expo-app@latest frontend

cd frontend

# 必要なパッケージのインストール
npm install @tanstack/react-query zustand axios react-hook-form zod
npm install nativewind tailwindcss
npm install @react-native-async-storage/async-storage
npm install @hookform/resolvers

# 開発サーバー起動
npx expo start

# iOS起動（シミュレータ）
npx expo start --ios

# Android起動（エミュレータ）
npx expo start --android

# Web版起動（開発時のみ）
npx expo start --web
```

**Note**: Expoを使用することで、物理デバイスでのテストも容易（Expo Goアプリ経由）。ただし、本格的なテストにはシミュレータ/エミュレータを推奨。

## テスト

### バックエンド
```bash
cd backend
go test ./... -v
go test ./... -cover
```

### フロントエンド
```bash
cd frontend
npm test
```

## トラブルシューティング

### データベース接続エラー
- Docker Composeが起動しているか確認: `docker-compose ps`
- 接続情報が正しいか確認: `.env`ファイル

### Expoビルドエラー
- キャッシュクリア: `npx expo start --clear`
- node_modules再インストール: `rm -rf node_modules && npm install`
- Metro Bundlerリセット: `npx expo start --reset-cache`
- iOS/Androidネイティブビルドが必要な場合: `npx expo prebuild`

## 参考リンク

### バックエンド
- [Chi Documentation](https://github.com/go-chi/chi)
- [sqlc Documentation](https://docs.sqlc.dev/)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [slog Documentation](https://pkg.go.dev/log/slog)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

### フロントエンド
- [Expo Documentation](https://docs.expo.dev/)
- [Expo Router Documentation](https://docs.expo.dev/router/introduction/)
- [NativeWind Documentation](https://www.nativewind.dev/)
- [TanStack Query Documentation](https://tanstack.com/query/latest/docs/react/overview)
- [Zustand Documentation](https://docs.pmnd.rs/zustand/)
- [React Hook Form](https://react-hook-form.com/)
- [Zod Documentation](https://zod.dev/)
