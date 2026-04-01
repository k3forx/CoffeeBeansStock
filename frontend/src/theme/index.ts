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

export const getAlertAccentColor = (
  alertLevel: "danger" | "warning" | undefined,
  remainingDays: number | undefined,
): string => {
  if (alertLevel === "danger") return colors.danger;
  if (alertLevel === "warning") return colors.warning;
  if (remainingDays != null) return colors.success;
  return colors.accent;
};
