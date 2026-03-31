import React from "react";
import { render, fireEvent } from "@testing-library/react-native";
import { Button } from "../Button";

describe("Button", () => {
  const defaultProps = {
    title: "テスト",
    onPress: jest.fn(),
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("タイトルテキストを表示する", () => {
    const { getByText } = render(<Button {...defaultProps} />);
    expect(getByText("テスト")).toBeTruthy();
  });

  it("押下時にonPressが呼ばれる", () => {
    const onPress = jest.fn();
    const { getByText } = render(<Button {...defaultProps} onPress={onPress} />);
    fireEvent.press(getByText("テスト"));
    expect(onPress).toHaveBeenCalledTimes(1);
  });

  it("disabled=trueの場合に押下してもonPressが呼ばれない", () => {
    const onPress = jest.fn();
    const { getByText } = render(
      <Button {...defaultProps} onPress={onPress} disabled />
    );
    fireEvent.press(getByText("テスト"));
    expect(onPress).not.toHaveBeenCalled();
  });

  it("loading=trueの場合にloadingTextを表示する", () => {
    const { getByText, queryByText } = render(
      <Button {...defaultProps} loading loadingText="読込中..." />
    );
    expect(getByText("読込中...")).toBeTruthy();
    expect(queryByText("テスト")).toBeNull();
  });

  it("loading=trueの場合に押下してもonPressが呼ばれない", () => {
    const onPress = jest.fn();
    const { getByText } = render(
      <Button {...defaultProps} onPress={onPress} loading loadingText="読込中..." />
    );
    fireEvent.press(getByText("読込中..."));
    expect(onPress).not.toHaveBeenCalled();
  });

  it("loadingTextが未指定の場合はloading中もtitleを表示する", () => {
    const { getByText } = render(<Button {...defaultProps} loading />);
    expect(getByText("テスト")).toBeTruthy();
  });
});
