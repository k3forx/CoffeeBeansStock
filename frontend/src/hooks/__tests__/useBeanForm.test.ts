import { renderHook, act } from "@testing-library/react-native";
import { useBeanForm } from "../useBeanForm";
import type { CoffeeBean } from "@/types/api";

describe("useBeanForm", () => {
  it("initial state has empty fields and no errors", () => {
    const { result } = renderHook(() => useBeanForm());

    expect(result.current.fields).toEqual({
      name: "",
      origin: "",
      roastLevel: "",
      roastDetail: "",
      currentStock: "",
      notes: "",
    });
    expect(result.current.errors).toEqual({});
  });

  it("setField updates the specified field", () => {
    const { result } = renderHook(() => useBeanForm());

    act(() => {
      result.current.setField("name", "Ethiopia Yirgacheffe");
    });

    expect(result.current.fields.name).toBe("Ethiopia Yirgacheffe");
  });

  it("setField clears the error for that field", () => {
    const { result } = renderHook(() => useBeanForm());

    // Trigger validation to produce errors
    act(() => {
      result.current.validate();
    });
    expect(result.current.errors.name).toBeDefined();

    // Setting the field should clear its error
    act(() => {
      result.current.setField("name", "Test Bean");
    });
    expect(result.current.errors.name).toBeUndefined();
    // Other errors should remain
    expect(result.current.errors.roastLevel).toBeDefined();
  });

  it("handleRoastLevelSelect sets roastLevel and clears roastDetail", () => {
    const { result } = renderHook(() => useBeanForm());

    // First set a detail
    act(() => {
      result.current.setField("roastDetail", "french");
    });
    expect(result.current.fields.roastDetail).toBe("french");

    // Selecting a new roast level should clear detail
    act(() => {
      result.current.handleRoastLevelSelect("medium");
    });
    expect(result.current.fields.roastLevel).toBe("medium");
    expect(result.current.fields.roastDetail).toBe("");
  });

  it("validate() returns errors for empty required fields and sets errors state", () => {
    const { result } = renderHook(() => useBeanForm());

    let validationResult: ReturnType<typeof result.current.validate>;
    act(() => {
      validationResult = result.current.validate();
    });

    expect(validationResult!.valid).toBe(false);
    if (!validationResult!.valid) {
      expect(validationResult!.errors.name).toBeDefined();
      expect(validationResult!.errors.roastLevel).toBeDefined();
      expect(validationResult!.errors.currentStock).toBeDefined();
    }
    // Errors should also be set in state
    expect(result.current.errors.name).toBeDefined();
  });

  it("validate() returns valid:true with parsed stock for valid input", () => {
    const { result } = renderHook(() => useBeanForm());

    act(() => {
      result.current.setField("name", "Colombia");
      result.current.setField("roastLevel", "deep");
      result.current.setField("currentStock", "250");
    });

    let validationResult: ReturnType<typeof result.current.validate>;
    act(() => {
      validationResult = result.current.validate();
    });

    expect(validationResult!).toEqual({ valid: true, stock: 250 });
    expect(result.current.errors).toEqual({});
  });

  it("resetToBean sets all fields from a CoffeeBean object", () => {
    const { result } = renderHook(() => useBeanForm());

    const bean: CoffeeBean = {
      id: "abc-123",
      name: "Guatemala Antigua",
      origin: "Guatemala",
      roast_level: "medium_deep",
      roast_detail: "city",
      current_stock: 300,
      notes: "Chocolate notes",
      created_at: "2025-01-01T00:00:00Z",
      updated_at: "2025-01-02T00:00:00Z",
      consumption_rate: {
        remaining_cups: 20,
        remaining_days: 14,
        daily_consumption_grams: 15.0,
        weekly_total_grams: 105,
      },
    };

    act(() => {
      result.current.resetToBean(bean);
    });

    expect(result.current.fields).toEqual({
      name: "Guatemala Antigua",
      origin: "Guatemala",
      roastLevel: "medium_deep",
      roastDetail: "city",
      currentStock: "300",
      notes: "Chocolate notes",
    });
    expect(result.current.errors).toEqual({});
  });

  it("toCreateInput produces correct API payload", () => {
    const { result } = renderHook(() => useBeanForm());

    act(() => {
      result.current.setField("name", "Kenya AA");
      result.current.setField("origin", "Kenya");
      result.current.setField("roastLevel", "shallow");
      result.current.setField("roastDetail", "cinnamon");
      result.current.setField("currentStock", "150");
      result.current.setField("notes", "Bright acidity");
    });

    const input = result.current.toCreateInput(150);

    expect(input).toEqual({
      name: "Kenya AA",
      origin: "Kenya",
      roast_level: "shallow",
      roast_detail: "cinnamon",
      current_stock: 150,
      notes: "Bright acidity",
    });
  });

  it("toCreateInput omits empty optional fields", () => {
    const { result } = renderHook(() => useBeanForm());

    act(() => {
      result.current.setField("name", "Brazil");
      result.current.setField("roastLevel", "medium");
      result.current.setField("currentStock", "100");
    });

    const input = result.current.toCreateInput(100);

    expect(input).toEqual({
      name: "Brazil",
      origin: undefined,
      roast_level: "medium",
      roast_detail: undefined,
      current_stock: 100,
      notes: undefined,
    });
  });

  it("toUpdateInput produces correct API payload", () => {
    const { result } = renderHook(() => useBeanForm());

    act(() => {
      result.current.setField("name", "Ethiopia Sidamo");
      result.current.setField("origin", "Ethiopia");
      result.current.setField("roastLevel", "deep");
      result.current.setField("roastDetail", "french");
      result.current.setField("notes", "Smoky");
    });

    const input = result.current.toUpdateInput(200);

    expect(input).toEqual({
      name: "Ethiopia Sidamo",
      origin: "Ethiopia",
      roast_level: "deep",
      roast_detail: "french",
      current_stock: 200,
      notes: "Smoky",
    });
  });
});
