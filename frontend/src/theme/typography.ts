import { Platform, TextStyle } from "react-native";

const serifFont = Platform.select({
  ios: "Hiragino Mincho ProN",
  android: "Noto Serif CJK JP",
  default: "serif",
});

const sansFont = Platform.select({
  ios: "Hiragino Sans",
  android: "Roboto",
  default: "System",
});

export const fontFamily = { serif: serifFont, sans: sansFont } as const;

export const typography = {
  displayLarge: {
    fontFamily: serifFont,
    fontSize: 20,
    fontWeight: "500" as TextStyle["fontWeight"],
    lineHeight: 28,
    letterSpacing: 3,
  },
  displaySub: {
    fontFamily: "Fraunces" as string,
    fontSize: 14,
    fontWeight: "300" as TextStyle["fontWeight"],
    fontStyle: "italic" as TextStyle["fontStyle"],
    lineHeight: 20,
    letterSpacing: 2,
  },
  titleLarge: {
    fontFamily: serifFont,
    fontSize: 22,
    fontWeight: "500" as TextStyle["fontWeight"],
    lineHeight: 28,
    letterSpacing: 1,
  },
  titleMedium: {
    fontFamily: serifFont,
    fontSize: 16,
    fontWeight: "500" as TextStyle["fontWeight"],
    lineHeight: 22,
    letterSpacing: 0.5,
  },
  bodyLarge: {
    fontFamily: sansFont,
    fontSize: 16,
    fontWeight: "400" as TextStyle["fontWeight"],
    lineHeight: 24,
  },
  bodyMedium: {
    fontFamily: sansFont,
    fontSize: 14,
    fontWeight: "400" as TextStyle["fontWeight"],
    lineHeight: 20,
  },
  bodySmall: {
    fontFamily: sansFont,
    fontSize: 12,
    fontWeight: "400" as TextStyle["fontWeight"],
    lineHeight: 18,
    letterSpacing: 0.3,
  },
  labelLarge: {
    fontFamily: sansFont,
    fontSize: 14,
    fontWeight: "600" as TextStyle["fontWeight"],
    lineHeight: 20,
    letterSpacing: 0.5,
  },
  labelMedium: {
    fontFamily: sansFont,
    fontSize: 12,
    fontWeight: "600" as TextStyle["fontWeight"],
    lineHeight: 16,
    letterSpacing: 0.5,
  },
  labelSmall: {
    fontFamily: sansFont,
    fontSize: 10,
    fontWeight: "500" as TextStyle["fontWeight"],
    lineHeight: 14,
    letterSpacing: 0.5,
  },
  stockNumber: {
    fontFamily: serifFont,
    fontSize: 16,
    fontWeight: "400" as TextStyle["fontWeight"],
    lineHeight: 18,
  },
  stockNumberLarge: {
    fontFamily: serifFont,
    fontSize: 40,
    fontWeight: "300" as TextStyle["fontWeight"],
    lineHeight: 44,
  },
} as const;
