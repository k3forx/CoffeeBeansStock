export { colors } from "./colors";
export { typography, fontFamily } from "./typography";
export { spacing, radius, shadows } from "./spacing";
export { formStyles } from "./formStyles";

import { colors } from "./colors";

export const getStockColor = (stock: number) => {
  if (stock <= 50) return colors.danger;
  if (stock <= 100) return colors.warning;
  return colors.success;
};

export const getStockAlertLevel = (
  remainingDays: number | null | undefined,
): "danger" | "warning" | null => {
  if (remainingDays == null) return null;
  if (remainingDays <= 3) return "danger";
  if (remainingDays <= 7) return "warning";
  return null;
};

export const getAlertAccentColor = (
  remainingDays: number | null | undefined,
): string => {
  if (remainingDays == null) return colors.accent;
  if (remainingDays <= 3) return colors.danger;
  if (remainingDays <= 7) return colors.warning;
  return colors.success;
};
