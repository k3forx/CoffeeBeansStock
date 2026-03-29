import type { RoastLevel, RoastDetail } from "../types/api";

export const ROAST_LEVELS: { value: RoastLevel; label: string }[] = [
  { value: "shallow", label: "浅煎り" },
  { value: "medium", label: "中煎り" },
  { value: "medium_deep", label: "中深煎り" },
  { value: "deep", label: "深煎り" },
];

export const ROAST_DETAILS: Record<RoastLevel, { value: RoastDetail; label: string }[]> = {
  shallow: [
    { value: "light", label: "ライトロースト" },
    { value: "cinnamon", label: "シナモンロースト" },
  ],
  medium: [
    { value: "medium", label: "ミディアムロースト" },
    { value: "high", label: "ハイロースト" },
  ],
  medium_deep: [
    { value: "city", label: "シティロースト" },
    { value: "full_city", label: "フルシティロースト" },
  ],
  deep: [
    { value: "french", label: "フレンチロースト" },
    { value: "italian", label: "イタリアンロースト" },
  ],
};

export const ROAST_LEVEL_LABELS: Record<RoastLevel, string> = {
  shallow: "浅煎り",
  medium: "中煎り",
  medium_deep: "中深煎り",
  deep: "深煎り",
};

export const ROAST_DETAIL_LABELS: Record<RoastDetail, string> = {
  light: "ライトロースト",
  cinnamon: "シナモンロースト",
  medium: "ミディアムロースト",
  high: "ハイロースト",
  city: "シティロースト",
  full_city: "フルシティロースト",
  french: "フレンチロースト",
  italian: "イタリアンロースト",
};
