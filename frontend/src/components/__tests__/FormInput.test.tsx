import React from "react";
import { render, fireEvent } from "@testing-library/react-native";
import { FormInput } from "../FormInput";

describe("FormInput", () => {
  const defaultProps = {
    label: "豆の名前",
    value: "",
    onChangeText: jest.fn(),
  };

  it("ラベルテキストを表示する", () => {
    const { getByText } = render(<FormInput {...defaultProps} />);
    expect(getByText("豆の名前")).toBeTruthy();
  });

  it("required=trueの場合に必須マークを表示する", () => {
    const { getByText } = render(
      <FormInput {...defaultProps} required={true} />
    );
    expect(getByText(" *")).toBeTruthy();
  });

  it("required=falseの場合に必須マークを表示しない", () => {
    const { queryByText } = render(
      <FormInput {...defaultProps} required={false} />
    );
    expect(queryByText(" *")).toBeNull();
  });

  it("入力変更時にonChangeTextが呼ばれる", () => {
    const onChangeText = jest.fn();
    const { getByDisplayValue } = render(
      <FormInput {...defaultProps} value="test" onChangeText={onChangeText} />
    );
    fireEvent.changeText(getByDisplayValue("test"), "new value");
    expect(onChangeText).toHaveBeenCalledWith("new value");
  });

  it("errorが指定されている場合にエラーメッセージを表示する", () => {
    const { getByText } = render(
      <FormInput {...defaultProps} error="入力必須です" />
    );
    expect(getByText("入力必須です")).toBeTruthy();
  });

  it("errorが未指定の場合にエラーメッセージを表示しない", () => {
    const { queryByText } = render(<FormInput {...defaultProps} />);
    expect(queryByText("入力必須です")).toBeNull();
  });
});
