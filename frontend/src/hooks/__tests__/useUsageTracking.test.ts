import { renderHook, act } from "@testing-library/react-native";
import { Alert } from "react-native";
import { useUsageTracking } from "@/hooks/useUsageTracking";
import { beansApi } from "@/api/beans";
import { usagesApi } from "@/api/usages";
import { showApiError } from "@/utils/errorHandler";
import type { CoffeeBean, UsageHistory } from "@/types/api";

jest.mock("@/api/beans");
jest.mock("@/api/usages");
jest.mock("@/utils/errorHandler");

const mockBeansGet = beansApi.get as jest.MockedFunction<typeof beansApi.get>;
const mockUsagesList = usagesApi.list as jest.MockedFunction<typeof usagesApi.list>;
const mockUsagesCreate = usagesApi.create as jest.MockedFunction<typeof usagesApi.create>;
const mockUsagesDelete = usagesApi.delete as jest.MockedFunction<typeof usagesApi.delete>;
const mockShowApiError = showApiError as jest.MockedFunction<typeof showApiError>;

jest.spyOn(Alert, "alert");

const fakeBean: CoffeeBean = {
  id: "bean-1",
  name: "Ethiopia",
  roast_level: "medium",
  current_stock: 200,
  created_at: "2025-01-01T00:00:00Z",
  updated_at: "2025-01-02T00:00:00Z",
  consumption_rate: {
    remaining_cups: 13,
  },
};

const fakeUsage: UsageHistory = {
  id: "usage-1",
  coffee_bean_id: "bean-1",
  usage_date: "2025-06-01",
  quantity: 15,
  created_at: "2025-06-01T00:00:00Z",
};

const fakePagination = { total: 1, limit: 10, offset: 0, has_more: false };

const defaultConfig = {
  beanId: "bean-1",
  onBeanUpdated: jest.fn(),
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe("useUsageTracking", () => {
  it("loadUsages fetches and sets usage list", async () => {
    mockUsagesList.mockResolvedValueOnce({
      usages: [fakeUsage],
      pagination: fakePagination,
    });

    const { result } = renderHook(() => useUsageTracking(defaultConfig));

    await act(async () => {
      await result.current.loadUsages();
    });

    expect(mockUsagesList).toHaveBeenCalledWith("bean-1", 10, 0);
    expect(result.current.usages).toEqual([fakeUsage]);
    expect(result.current.usageLoading).toBe(false);
  });

  it("handleQuickUsage calls usagesApi.create with correct params", async () => {
    mockUsagesCreate.mockResolvedValueOnce(fakeUsage);
    mockBeansGet.mockResolvedValueOnce(fakeBean);
    mockUsagesList.mockResolvedValueOnce({
      usages: [fakeUsage],
      pagination: fakePagination,
    });

    const { result } = renderHook(() => useUsageTracking(defaultConfig));

    await act(async () => {
      await result.current.handleQuickUsage(15);
    });

    expect(mockUsagesCreate).toHaveBeenCalledWith("bean-1", {
      usage_date: expect.stringMatching(/^\d{4}-\d{2}-\d{2}$/),
      quantity: 15,
    });
  });

  it("handleQuickUsage calls onBeanUpdated callback after success", async () => {
    const onBeanUpdated = jest.fn();
    mockUsagesCreate.mockResolvedValueOnce(fakeUsage);
    mockBeansGet.mockResolvedValueOnce(fakeBean);
    mockUsagesList.mockResolvedValueOnce({
      usages: [fakeUsage],
      pagination: fakePagination,
    });

    const { result } = renderHook(() =>
      useUsageTracking({ beanId: "bean-1", onBeanUpdated })
    );

    await act(async () => {
      await result.current.handleQuickUsage(15);
    });

    expect(onBeanUpdated).toHaveBeenCalledWith(fakeBean);
  });

  it("handleQuickUsage sets quickButtonLoading during request", async () => {
    let resolveCreate: (value: UsageHistory) => void;
    mockUsagesCreate.mockImplementationOnce(
      () => new Promise((resolve) => { resolveCreate = resolve; })
    );

    const { result } = renderHook(() => useUsageTracking(defaultConfig));

    // Start quick usage (don't await)
    let promise: Promise<void>;
    act(() => {
      promise = result.current.handleQuickUsage(15);
    });

    expect(result.current.quickButtonLoading).toBe(15);

    // Resolve and finish
    mockBeansGet.mockResolvedValueOnce(fakeBean);
    mockUsagesList.mockResolvedValueOnce({
      usages: [fakeUsage],
      pagination: fakePagination,
    });

    await act(async () => {
      resolveCreate!(fakeUsage);
      await promise;
    });

    expect(result.current.quickButtonLoading).toBe(null);
  });

  it("handleQuickUsage shows error on INSUFFICIENT_STOCK", async () => {
    const error = new Error("insufficient stock");
    mockUsagesCreate.mockRejectedValueOnce(error);

    const { result } = renderHook(() => useUsageTracking(defaultConfig));

    await act(async () => {
      await result.current.handleQuickUsage(15);
    });

    expect(mockShowApiError).toHaveBeenCalledWith(error, "記録に失敗しました", {
      INSUFFICIENT_STOCK: "在庫が不足しています",
    });
  });

  it("handleManualUsage validates grams > 0", async () => {
    const { result } = renderHook(() => useUsageTracking(defaultConfig));

    // Set invalid input
    act(() => {
      result.current.setManualGrams("0");
    });

    await act(async () => {
      await result.current.handleManualUsage();
    });

    expect(Alert.alert).toHaveBeenCalledWith("エラー", "1g以上の数値を入力してください");
    expect(mockUsagesCreate).not.toHaveBeenCalled();
  });

  it("handleManualUsage clears manualGrams on success", async () => {
    mockUsagesCreate.mockResolvedValueOnce(fakeUsage);
    mockBeansGet.mockResolvedValueOnce(fakeBean);
    mockUsagesList.mockResolvedValueOnce({
      usages: [fakeUsage],
      pagination: fakePagination,
    });

    const { result } = renderHook(() => useUsageTracking(defaultConfig));

    act(() => {
      result.current.setManualGrams("20");
    });

    await act(async () => {
      await result.current.handleManualUsage();
    });

    expect(result.current.manualGrams).toBe("");
  });

  it("handleDeleteUsage calls usagesApi.delete and refreshes data", async () => {
    mockUsagesDelete.mockResolvedValueOnce({ message: "deleted" });
    mockBeansGet.mockResolvedValueOnce(fakeBean);
    mockUsagesList.mockResolvedValueOnce({
      usages: [],
      pagination: fakePagination,
    });

    const { result } = renderHook(() => useUsageTracking(defaultConfig));

    await act(async () => {
      await result.current.handleDeleteUsage("usage-1");
    });

    expect(mockUsagesDelete).toHaveBeenCalledWith("bean-1", "usage-1");
    expect(mockBeansGet).toHaveBeenCalledWith("bean-1");
    expect(result.current.usages).toEqual([]);
  });

  it("handleDeleteUsage sets deletingUsageId during request", async () => {
    let resolveDelete: (value: { message: string }) => void;
    mockUsagesDelete.mockImplementationOnce(
      () => new Promise((resolve) => { resolveDelete = resolve; })
    );

    const { result } = renderHook(() => useUsageTracking(defaultConfig));

    let promise: Promise<void>;
    act(() => {
      promise = result.current.handleDeleteUsage("usage-1");
    });

    expect(result.current.deletingUsageId).toBe("usage-1");

    mockBeansGet.mockResolvedValueOnce(fakeBean);
    mockUsagesList.mockResolvedValueOnce({
      usages: [],
      pagination: fakePagination,
    });

    await act(async () => {
      resolveDelete!({ message: "deleted" });
      await promise;
    });

    expect(result.current.deletingUsageId).toBe(null);
  });
});
