# Phase 3: 画面コンポーネントテスト

## 目的

`@testing-library/react-native`を導入し、画面コンポーネントのインテグレーションテストを追加する。
ユーザー操作（タップ、入力）→画面表示の変化をテストし、UIのリグレッションを防ぐ。

## 前提条件

### react-test-renderer のバージョン競合

現在 `react@19.1.0` に対し、`@testing-library/react-native` が依存する `react-test-renderer` は `react@^19.2.4` を要求する。

**対応方法（いずれか）**:
1. **推奨**: `react` と `react-native` を最新にアップデートしてからインストール
2. **暫定**: `npm install --legacy-peer-deps` でインストール（互換性リスクあり）

### 追加パッケージ

```bash
npm install --save-dev @testing-library/react-native @testing-library/jest-native
```

### jest.setup.ts への追記

```typescript
import "@testing-library/jest-native/extend-expect";
```

## 対象画面と作成するテスト

### 1. `app/auth/__tests__/login.test.tsx`

**テスト対象**: `app/auth/login.tsx` — ログイン画面

**テストケース**:

| テスト | 検証内容 |
|--------|---------|
| 初期表示でログインボタンが表示される | `render` → ボタンテキストが存在するか |
| ボタン押下で匿名登録APIが呼ばれる | `fireEvent.press` → `authApi.registerAnonymous` が呼ばれるか |
| 登録成功後にauthストアが更新される | `setAuth` が正しい引数で呼ばれるか |
| エラー時にエラーメッセージが表示される | API失敗 → エラーテキストが画面に表示されるか |
| ローディング中はボタンが無効化される | ボタンの`disabled`状態 |

**モック**:
- `authApi` — `jest.mock("@/api/auth")`
- authストア — `useAuthStore.setState()` で初期状態設定

### 2. `app/(tabs)/__tests__/index.test.tsx`

**テスト対象**: `app/(tabs)/index.tsx` — 豆一覧画面

**テストケース**:

| テスト | 検証内容 |
|--------|---------|
| ローディング中にスピナーが表示される | `ActivityIndicator` が表示されるか |
| 豆一覧が正しく表示される | モックデータの豆名がすべて表示されるか |
| 空状態で適切なメッセージが表示される | 豆0件のときの表示 |
| 豆カードタップで詳細画面へ遷移する | `fireEvent.press` → `router.push` が `/beans/{id}` で呼ばれるか |
| 新規作成ボタンで作成画面へ遷移する | `router.push` が `/beans/create` で呼ばれるか |
| 在庫レベルに応じた色が表示される | `getStockColor` の結果がスタイルに反映されているか |

**モック**:
- `beansApi` — `jest.mock("@/api/beans")`
- `expo-router` — jest.setup.tsで設定済み（`mockRouter`をexportして検証に利用）

### 3. `app/(tabs)/beans/__tests__/create.test.tsx`

**テスト対象**: `app/(tabs)/beans/create.tsx` — 豆作成フォーム

**テストケース**:

| テスト | 検証内容 |
|--------|---------|
| フォームフィールドが表示される | 名前、在庫量、焙煎度のフィールドが存在するか |
| 名前未入力で送信するとバリデーションエラー | `Alert.alert` or エラーメッセージ表示 |
| 正常入力で送信するとAPIが呼ばれる | `beansApi.create` が正しい引数で呼ばれるか |
| 作成成功後に一覧画面へ戻る | `router.back()` or `router.replace()` が呼ばれるか |
| API失敗時にエラーが表示される | エラーメッセージの表示 |

**モック**:
- `beansApi` — `jest.mock("@/api/beans")`
- `Alert` — `jest.spyOn(Alert, "alert")`（バリデーションにAlertを使っている場合）

### 4. `app/(tabs)/beans/__tests__/[id].test.tsx`

**テスト対象**: `app/(tabs)/beans/[id].tsx` — 豆詳細画面

**テストケース**:

| テスト | 検証内容 |
|--------|---------|
| 豆の詳細情報が表示される | 名前、在庫量、焙煎度、消費ペースなどが正しく表示されるか |
| 編集モードに切り替わる | 編集ボタンタップ → フォームが表示されるか |
| 保存でAPIが呼ばれる | `beansApi.update` が呼ばれるか |
| 削除確認ダイアログが表示される | 削除ボタン → `Alert.alert` で確認 |
| 削除確認で実行するとAPIが呼ばれ一覧に戻る | `beansApi.delete` → `router.back()` |

**モック**:
- `beansApi` — `jest.mock("@/api/beans")`
- `useLocalSearchParams` — `{ id: "test-id" }` を返す
- `Alert` — `jest.spyOn(Alert, "alert")`

### 5. `app/(tabs)/__tests__/profile.test.tsx`

**テスト対象**: `app/(tabs)/profile.tsx` — プロフィール画面

**テストケース**:

| テスト | 検証内容 |
|--------|---------|
| ユーザー情報が表示される | ストアのuser情報が表示されるか |
| ログアウトボタンが動作する | タップ → 確認ダイアログ → `logout` 呼び出し |

### 6. `app/__tests__/_layout.test.tsx`

**テスト対象**: `app/_layout.tsx` — ルートレイアウト（認証ルーティング）

**テストケース**:

| テスト | 検証内容 |
|--------|---------|
| 未認証時にログイン画面へリダイレクト | `isAuthenticated=false` → `router.replace("/auth/login")` |
| 認証済みでタブ画面が表示される | `isAuthenticated=true` → リダイレクトしない |
| ローディング中はスプラッシュ/スピナー表示 | `isLoading=true` → ナビゲーションしない |

## jest.setup.ts の更新

```typescript
// expo-router のモックをexport可能にする
const mockRouter = {
  push: jest.fn(),
  replace: jest.fn(),
  back: jest.fn(),
};

jest.mock("expo-router", () => ({
  useRouter: () => mockRouter,
  useSegments: jest.fn(() => []),
  useLocalSearchParams: jest.fn(() => ({})),
  useFocusEffect: (cb: () => void) => cb(),
  Link: ({ children }: { children: React.ReactNode }) => children,
  Slot: () => null,
}));

export { mockRouter };
```

## 実装手順

1. React / react-native のバージョン互換性を確認・必要ならアップデート
2. `@testing-library/react-native` をインストール
3. `jest.setup.ts` を更新（mockRouter export、extend-expect追加）
4. `app/auth/__tests__/login.test.tsx` を作成（最もシンプルな画面から開始）
5. `app/(tabs)/__tests__/index.test.tsx` を作成
6. `app/(tabs)/beans/__tests__/create.test.tsx` を作成
7. `app/(tabs)/beans/__tests__/[id].test.tsx` を作成
8. `app/(tabs)/__tests__/profile.test.tsx` を作成
9. `app/__tests__/_layout.test.tsx` を作成
10. カバレッジ閾値を設定（`package.json` の jest設定に追加）

```json
"coverageThreshold": {
  "global": {
    "branches": 60,
    "functions": 70,
    "lines": 70,
    "statements": 70
  }
}
```

## 検証コマンド

```bash
cd frontend
npm test -- --verbose
npm test -- --coverage
```
