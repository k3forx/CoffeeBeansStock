---
paths:
  - "frontend/src/**/*.test.ts"
  - "frontend/src/**/*.test.tsx"
---

# フロントエンド テスト方針

## テストレイヤーと対象

| レイヤー | 対象 | テストツール | 方針 |
|----------|------|-------------|------|
| **Unit** | utils, constants, 純粋関数 | Jest | 入出力のみ検証 |
| **Hook** | custom hooks | `renderHook` + `act` + mock | API層をモック、状態遷移を検証 |
| **Store** | Zustand store | `renderHook` or 直接呼び出し | SecureStoreはjest.setupでモック済み、状態変更を検証 |
| **Component** | 再利用コンポーネント | `render` (@testing-library/react-native) | props→表示・イベント→コールバック呼び出しを検証 |

### テスト対象外

- **API wrappers** (`api/beans.ts`, `api/auth.ts` 等): `api/client.ts`への薄い委譲のみでロジックなし。`client.test.ts`でカバー済み
- **画面コンポーネント** (`app/` 配下): モック量が多くメンテコスト高。hooks層テストで間接カバー。必要時に段階的追加
- **型定義** (`types/`): 実行時ロジックなし
- **定数・テーマ** (`theme/colors.ts`, `constants/` 等): 定数定義のみ

## テストファイル配置

- ソースと同階層の `__tests__/` ディレクトリに配置
- 命名規則: `{対象ファイル名}.test.ts` (`.tsx` はJSXを含む場合)

```
src/
├── hooks/
│   ├── __tests__/
│   │   └── useBeansList.test.ts
│   └── useBeansList.ts
├── components/
│   ├── __tests__/
│   │   └── FormInput.test.tsx
│   └── FormInput.tsx
```

## テストの書き方ルール

### 共通

- テスト名は日本語で記述
- `beforeEach(() => jest.clearAllMocks())` でモック状態をリセット
- 外部依存（API, SecureStore等）のみモック。内部ロジックはモックしない

### Hook テスト テンプレート

```typescript
import { renderHook, act } from "@testing-library/react-native";
import { useXxx } from "@/hooks/useXxx";
import { xxxApi } from "@/api/xxx";
import { showApiError } from "@/utils/errorHandler";

jest.mock("@/api/xxx");
jest.mock("@/utils/errorHandler");

const mockMethod = xxxApi.method as jest.MockedFunction<typeof xxxApi.method>;

beforeEach(() => {
  jest.clearAllMocks();
});

describe("useXxx", () => {
  it("呼び出し時にXxxを実行する", async () => {
    mockMethod.mockResolvedValueOnce(fakeData);

    const { result } = renderHook(() => useXxx());

    await act(async () => {
      await result.current.someAction();
    });

    expect(result.current.state).toEqual(expected);
  });
});
```

### Component テスト テンプレート

```tsx
import { render, fireEvent } from "@testing-library/react-native";
import { XxxComponent } from "@/components/XxxComponent";

describe("XxxComponent", () => {
  const defaultProps = {
    // required props with sensible defaults
  };

  it("ラベルテキストを表示する", () => {
    const { getByText } = render(<XxxComponent {...defaultProps} />);
    expect(getByText("expected text")).toBeTruthy();
  });

  it("押下時にonXxxが呼ばれる", () => {
    const onXxx = jest.fn();
    const { getByText } = render(
      <XxxComponent {...defaultProps} onXxx={onXxx} />
    );

    fireEvent.press(getByText("button text"));
    expect(onXxx).toHaveBeenCalledWith(expectedArg);
  });
});
```

### Store テスト テンプレート

```typescript
import { useAuthStore } from "@/stores/auth";
import * as SecureStore from "expo-secure-store";

// jest.setupでSecureStoreはモック済み

beforeEach(() => {
  jest.clearAllMocks();
  // Zustand storeをリセット
  useAuthStore.setState({
    user: null,
    accessToken: null,
    refreshToken: null,
    isLoading: true,
    isAuthenticated: false,
  });
});

describe("useAuthStore", () => {
  it("setAuthでトークンを保存し状態を更新する", async () => {
    await useAuthStore.getState().setAuth(fakeUser, "access", "refresh");

    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(true);
    expect(SecureStore.setItemAsync).toHaveBeenCalledWith("access_token", "access");
  });
});
```

## アサーション方針

- 状態オブジェクト全体の比較を優先（`toEqual`）
- フィールド個別の `toBe` は状態遷移の検証（loading, refreshing等）に限定
- コンポーネントテストでは `getByText` / `getByTestId` で要素の存在確認、`fireEvent` でインタラクション検証
- スナップショットテストは使わない（壊れやすくレビューしにくい）
