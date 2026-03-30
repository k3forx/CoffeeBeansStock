import { Alert } from "react-native";
import { ApiError } from "@/api/client";
import { showApiError } from "@/utils/errorHandler";

jest.spyOn(Alert, "alert");

beforeEach(() => {
  jest.clearAllMocks();
});

describe("showApiError", () => {
  it("shows overridden message when ApiError code matches codeMessages", () => {
    const error = new ApiError("NOT_FOUND", "豆が見つかりません");

    showApiError(error, "エラーが発生しました", {
      NOT_FOUND: "指定された豆は存在しません",
    });

    expect(Alert.alert).toHaveBeenCalledWith(
      "エラー",
      "指定された豆は存在しません",
    );
  });

  it("shows error.message when ApiError code has no override", () => {
    const error = new ApiError("VALIDATION", "バリデーションエラー");

    showApiError(error);

    expect(Alert.alert).toHaveBeenCalledWith("エラー", "バリデーションエラー");
  });

  it("shows fallback message for non-ApiError", () => {
    const error = new TypeError("network error");

    showApiError(error);

    expect(Alert.alert).toHaveBeenCalledWith("エラー", "エラーが発生しました");
  });

  it("shows custom fallback message for non-ApiError", () => {
    const error = new Error("something went wrong");

    showApiError(error, "カスタムエラーメッセージ");

    expect(Alert.alert).toHaveBeenCalledWith(
      "エラー",
      "カスタムエラーメッセージ",
    );
  });
});
