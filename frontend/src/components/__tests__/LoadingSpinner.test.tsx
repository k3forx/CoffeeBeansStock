import React from "react";
import { render } from "@testing-library/react-native";
import { LoadingSpinner } from "../LoadingSpinner";

describe("LoadingSpinner", () => {
  it("デフォルトでlargeサイズのActivityIndicatorを表示する", () => {
    const { getByTestId } = render(<LoadingSpinner />);
    const spinner = getByTestId("loading-spinner");
    expect(spinner).toBeTruthy();
    expect(spinner.props.size).toBe("large");
  });

  it("sizeプロパティでサイズを変更できる", () => {
    const { getByTestId } = render(<LoadingSpinner size="small" />);
    const spinner = getByTestId("loading-spinner");
    expect(spinner.props.size).toBe("small");
  });
});
