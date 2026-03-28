export { colors } from "./colors";
export { typography, fontFamily } from "./typography";
export { spacing, radius, shadows } from "./spacing";

import { colors } from "./colors";

export const getStockColor = (stock: number) => {
  if (stock <= 50) return colors.danger;
  if (stock <= 100) return colors.warning;
  return colors.success;
};
