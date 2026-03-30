import { validateBeanForm } from "@/utils/validation";

describe("validateBeanForm", () => {
  it("returns error when name is empty", () => {
    const result = validateBeanForm({
      name: "",
      roastLevel: "Medium",
      currentStock: "100",
    });

    expect(result).toEqual({
      valid: false,
      errors: { name: "名前は必須です" },
    });
  });

  it("returns error when roastLevel is empty", () => {
    const result = validateBeanForm({
      name: "Ethiopia",
      roastLevel: "",
      currentStock: "100",
    });

    expect(result).toEqual({
      valid: false,
      errors: { roastLevel: "焙煎度を選択してください" },
    });
  });

  it("returns error when currentStock is non-numeric", () => {
    const result = validateBeanForm({
      name: "Ethiopia",
      roastLevel: "Medium",
      currentStock: "abc",
    });

    expect(result).toEqual({
      valid: false,
      errors: { currentStock: "在庫数は0以上の数値を入力してください" },
    });
  });

  it("returns error when currentStock is negative", () => {
    const result = validateBeanForm({
      name: "Ethiopia",
      roastLevel: "Medium",
      currentStock: "-1",
    });

    expect(result).toEqual({
      valid: false,
      errors: { currentStock: "在庫数は0以上の数値を入力してください" },
    });
  });

  it("returns all errors when all fields are invalid", () => {
    const result = validateBeanForm({
      name: "",
      roastLevel: "",
      currentStock: "abc",
    });

    expect(result).toEqual({
      valid: false,
      errors: {
        name: "名前は必須です",
        roastLevel: "焙煎度を選択してください",
        currentStock: "在庫数は0以上の数値を入力してください",
      },
    });
  });

  it("returns valid result with parsed stock for valid input", () => {
    const result = validateBeanForm({
      name: "Ethiopia",
      roastLevel: "Medium",
      currentStock: "250",
    });

    expect(result).toEqual({ valid: true, stock: 250 });
  });

  it("accepts stock of 0 as valid", () => {
    const result = validateBeanForm({
      name: "Ethiopia",
      roastLevel: "Medium",
      currentStock: "0",
    });

    expect(result).toEqual({ valid: true, stock: 0 });
  });
});
