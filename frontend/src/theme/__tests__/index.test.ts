import { getStockColor } from "../index";
import { colors } from "../colors";

describe("getStockColor", () => {
  it("returns danger for stock <= 50", () => {
    expect(getStockColor(0)).toBe(colors.danger);
    expect(getStockColor(50)).toBe(colors.danger);
  });

  it("returns warning for stock 51-100", () => {
    expect(getStockColor(51)).toBe(colors.warning);
    expect(getStockColor(100)).toBe(colors.warning);
  });

  it("returns success for stock > 100", () => {
    expect(getStockColor(101)).toBe(colors.success);
    expect(getStockColor(500)).toBe(colors.success);
  });
});
