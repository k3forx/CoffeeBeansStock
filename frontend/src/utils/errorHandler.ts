import { Alert } from "react-native";
import { ApiError } from "../api/client";

export function showApiError(
  e: unknown,
  fallback = "エラーが発生しました",
  codeMessages?: Record<string, string>,
): void {
  if (e instanceof ApiError) {
    const overridden = codeMessages?.[e.code];
    Alert.alert("エラー", overridden ?? e.message);
  } else {
    Alert.alert("エラー", fallback);
  }
}
