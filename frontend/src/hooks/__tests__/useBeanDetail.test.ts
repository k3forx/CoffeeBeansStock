import { renderHook, act } from "@testing-library/react-native";
import { useBeanDetail } from "@/hooks/useBeanDetail";
import { beansApi } from "@/api/beans";
import { showApiError } from "@/utils/errorHandler";
import type { CoffeeBean } from "@/types/api";

jest.mock("@/api/beans");
jest.mock("@/utils/errorHandler");

const mockGet = beansApi.get as jest.MockedFunction<typeof beansApi.get>;
const mockUpdate = beansApi.update as jest.MockedFunction<typeof beansApi.update>;
const mockDelete = beansApi.delete as jest.MockedFunction<typeof beansApi.delete>;
const mockShowApiError = showApiError as jest.MockedFunction<typeof showApiError>;

const fakeBean: CoffeeBean = {
  id: "bean-1",
  name: "Ethiopia Yirgacheffe",
  origin: "Ethiopia",
  roast_level: "medium",
  roast_detail: "high",
  current_stock: 200,
  notes: "Fruity",
  created_at: "2025-01-01T00:00:00Z",
  updated_at: "2025-01-02T00:00:00Z",
  consumption_rate: {
    remaining_cups: 13,
    remaining_days: 10,
    daily_consumption_grams: 15.0,
    weekly_total_grams: 105,
  },
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe("useBeanDetail", () => {
  it("loadBean calls beansApi.get and sets bean state", async () => {
    mockGet.mockResolvedValueOnce(fakeBean);

    const { result } = renderHook(() => useBeanDetail());

    await act(async () => {
      await result.current.loadBean("bean-1");
    });

    expect(mockGet).toHaveBeenCalledWith("bean-1");
    expect(result.current.bean).toEqual(fakeBean);
    expect(result.current.loading).toBe(false);
  });

  it("loadBean sets loading=false on error", async () => {
    const error = new Error("network error");
    mockGet.mockRejectedValueOnce(error);

    const { result } = renderHook(() => useBeanDetail());

    await act(async () => {
      await result.current.loadBean("bean-1");
    });

    expect(result.current.loading).toBe(false);
    expect(mockShowApiError).toHaveBeenCalledWith(error, "データの取得に失敗しました");
  });

  it("updateBean calls beansApi.update and sets bean, editing=false", async () => {
    const updatedBean = { ...fakeBean, name: "Updated Bean" };
    mockUpdate.mockResolvedValueOnce(updatedBean);

    const { result } = renderHook(() => useBeanDetail());

    // Start editing first
    act(() => {
      result.current.startEditing();
    });
    expect(result.current.editing).toBe(true);

    const input = { name: "Updated Bean" };
    await act(async () => {
      await result.current.updateBean("bean-1", input);
    });

    expect(mockUpdate).toHaveBeenCalledWith("bean-1", input);
    expect(result.current.bean).toEqual(updatedBean);
    expect(result.current.editing).toBe(false);
    expect(result.current.saving).toBe(false);
  });

  it("updateBean keeps editing=true on API error", async () => {
    const error = new Error("update failed");
    mockUpdate.mockRejectedValueOnce(error);

    const { result } = renderHook(() => useBeanDetail());

    act(() => {
      result.current.startEditing();
    });

    await act(async () => {
      await result.current.updateBean("bean-1", { name: "Fail" });
    });

    expect(result.current.editing).toBe(true);
    expect(result.current.saving).toBe(false);
    expect(mockShowApiError).toHaveBeenCalledWith(error, "更新に失敗しました");
  });

  it("deleteBean calls beansApi.delete and returns true on success", async () => {
    mockDelete.mockResolvedValueOnce({ message: "deleted" });

    const { result } = renderHook(() => useBeanDetail());

    let success: boolean;
    await act(async () => {
      success = await result.current.deleteBean("bean-1");
    });

    expect(mockDelete).toHaveBeenCalledWith("bean-1");
    expect(success!).toBe(true);
  });

  it("deleteBean returns false on error", async () => {
    const error = new Error("delete failed");
    mockDelete.mockRejectedValueOnce(error);

    const { result } = renderHook(() => useBeanDetail());

    let success: boolean;
    await act(async () => {
      success = await result.current.deleteBean("bean-1");
    });

    expect(success!).toBe(false);
    expect(mockShowApiError).toHaveBeenCalledWith(error, "削除に失敗しました");
  });

  it("startEditing/cancelEditing toggle editing state", () => {
    const { result } = renderHook(() => useBeanDetail());

    expect(result.current.editing).toBe(false);

    act(() => {
      result.current.startEditing();
    });
    expect(result.current.editing).toBe(true);

    act(() => {
      result.current.cancelEditing();
    });
    expect(result.current.editing).toBe(false);
  });
});
