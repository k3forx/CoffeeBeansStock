import { renderHook, act } from "@testing-library/react-native";
import { useBeansList } from "@/hooks/useBeansList";
import { beansApi } from "@/api/beans";
import { showApiError } from "@/utils/errorHandler";
import type { ListBeansResult } from "@/types/api";

jest.mock("@/api/beans");
jest.mock("@/utils/errorHandler");

const mockList = beansApi.list as jest.MockedFunction<typeof beansApi.list>;
const mockShowApiError = showApiError as jest.MockedFunction<typeof showApiError>;

const fakeBeans = [
  {
    id: "1",
    name: "Ethiopia",
    roast_level: "medium" as const,
    current_stock: 200,
    created_at: "2025-01-01T00:00:00Z",
    updated_at: "2025-01-01T00:00:00Z",
    consumption_rate: {
      remaining_cups: 13,
    },
  },
];

beforeEach(() => {
  jest.clearAllMocks();
});

describe("useBeansList", () => {
  it("fetchBeans calls beansApi.list and sets beans", async () => {
    mockList.mockResolvedValueOnce({
      beans: fakeBeans,
      pagination: { total: 1, limit: 100, offset: 0, has_more: false },
    });

    const { result } = renderHook(() => useBeansList());

    await act(async () => {
      await result.current.fetchBeans();
    });

    expect(mockList).toHaveBeenCalledWith(100, 0);
    expect(result.current.beans).toEqual(fakeBeans);
  });

  it("fetchBeans sets loading=false after successful fetch", async () => {
    mockList.mockResolvedValueOnce({
      beans: fakeBeans,
      pagination: { total: 1, limit: 100, offset: 0, has_more: false },
    });

    const { result } = renderHook(() => useBeansList());

    expect(result.current.loading).toBe(true);

    await act(async () => {
      await result.current.fetchBeans();
    });

    expect(result.current.loading).toBe(false);
  });

  it("fetchBeans sets loading=false on error and calls showApiError", async () => {
    const error = new Error("network error");
    mockList.mockRejectedValueOnce(error);

    const { result } = renderHook(() => useBeansList());

    await act(async () => {
      await result.current.fetchBeans();
    });

    expect(result.current.loading).toBe(false);
    expect(mockShowApiError).toHaveBeenCalledWith(error, "データの取得に失敗しました");
  });

  it("onRefresh sets refreshing=true then back to false", async () => {
    let resolveList: (value: ListBeansResult) => void;
    mockList.mockImplementationOnce(
      () => new Promise((resolve) => { resolveList = resolve; })
    );

    const { result } = renderHook(() => useBeansList());

    // Start refresh
    act(() => {
      result.current.onRefresh();
    });

    expect(result.current.refreshing).toBe(true);

    // Resolve the API call
    await act(async () => {
      resolveList!({
        beans: fakeBeans,
        pagination: { total: 1, limit: 100, offset: 0, has_more: false },
      });
    });

    expect(result.current.refreshing).toBe(false);
  });

  it("onRefresh calls beansApi.list to re-fetch", async () => {
    mockList.mockResolvedValueOnce({
      beans: fakeBeans,
      pagination: { total: 1, limit: 100, offset: 0, has_more: false },
    });

    const { result } = renderHook(() => useBeansList());

    await act(async () => {
      result.current.onRefresh();
    });

    expect(mockList).toHaveBeenCalledWith(100, 0);
  });
});
