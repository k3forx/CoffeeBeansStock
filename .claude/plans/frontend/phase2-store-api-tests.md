# Phase 2: ストア + API層テスト

## 目的

Zustand authストアとAPI薄ラッパー（auth, beans）のユニットテストを追加する。
トークン永続化やAPIのURL構築ミスはアプリ全体に影響するため、Phase 1の次に優先度が高い。

## 対象ファイルと作成するテスト

### 1. `src/stores/__tests__/auth.test.ts`

**テスト対象**: `src/stores/auth.ts` — Zustand authストア

**モック**: `expo-secure-store`（jest.setup.tsで設定済み）

**テストケース**:

| テスト | 検証内容 |
|--------|---------|
| `setAuth` がユーザー情報とトークンを保存 | `setState`後に`user`, `accessToken`, `refreshToken`, `isAuthenticated`が正しいか。`SecureStore.setItemAsync`が`access_token`, `refresh_token`で呼ばれるか |
| `setTokens` がトークンのみ更新 | `user`は変わらず、トークンだけ更新されるか |
| `logout` が全状態をクリア | `user=null`, `isAuthenticated=false`。`SecureStore.deleteItemAsync`が呼ばれるか |
| `loadStoredTokens` がトークンを復元 | SecureStoreにトークンがある場合→`isAuthenticated=true`、ない場合→`false` |
| `loadStoredTokens` が完了後に `isLoading=false` | 成功・失敗いずれでも`isLoading`が`false`になるか |

**実装メモ**:
```typescript
import * as SecureStore from "expo-secure-store";
import { useAuthStore } from "../auth";

// SecureStoreのモックは jest.setup.ts で設定済み
// テストごとにストアをリセット:
beforeEach(() => {
  useAuthStore.setState({
    user: null, accessToken: null, refreshToken: null,
    isAuthenticated: false, isLoading: true,
  });
  jest.clearAllMocks();
});

// loadStoredTokens のテストでは getItemAsync のモック戻り値を変更:
(SecureStore.getItemAsync as jest.Mock).mockResolvedValueOnce("stored-access");
```

### 2. `src/api/__tests__/auth.test.ts`

**テスト対象**: `src/api/auth.ts` — 認証API

**モック**: `global.fetch`（client.test.tsと同パターン）

**テストケース**:

| テスト | 検証内容 |
|--------|---------|
| `registerAnonymous` がPOST `/auth/register` を呼ぶ | URLとメソッドが正しいか |
| `getMe` がGET `/auth/me` を呼ぶ | URLとメソッドが正しいか |

### 3. `src/api/__tests__/beans.test.ts`

**テスト対象**: `src/api/beans.ts` — 豆CRUD API

**テストケース**:

| テスト | 検証内容 |
|--------|---------|
| `list` がGET `/coffee-beans/?limit=20&offset=0` を呼ぶ | デフォルト引数でのURL構築 |
| `list` にカスタムlimit/offsetを渡せる | `limit=10&offset=5`のURL構築 |
| `get` がGET `/coffee-beans/{id}` を呼ぶ | IDのパス埋め込み |
| `create` がPOST `/coffee-beans/` + bodyを呼ぶ | メソッドとbody |
| `update` がPUT `/coffee-beans/{id}` + bodyを呼ぶ | メソッドとbody |
| `delete` がDELETE `/coffee-beans/{id}` を呼ぶ | メソッド |

## 実装手順

1. `src/stores/__tests__/auth.test.ts` を作成・テスト実行
2. `src/api/__tests__/auth.test.ts` を作成・テスト実行
3. `src/api/__tests__/beans.test.ts` を作成・テスト実行
4. `npm test` で全テスト（Phase 1 含む）がパスすることを確認

## 検証コマンド

```bash
cd frontend
npm test -- --verbose
```
