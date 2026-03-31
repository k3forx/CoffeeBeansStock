import React from "react";
import { render, fireEvent } from "@testing-library/react-native";
import { ChipSelector } from "../ChipSelector";

const items = [
  { label: "浅煎り", value: "light" },
  { label: "中煎り", value: "medium" },
  { label: "深煎り", value: "dark" },
];

describe("ChipSelector", () => {
  it("全てのチップアイテムを表示する", () => {
    const { getByText } = render(
      <ChipSelector items={items} selected="" onSelect={jest.fn()} />
    );
    expect(getByText("浅煎り")).toBeTruthy();
    expect(getByText("中煎り")).toBeTruthy();
    expect(getByText("深煎り")).toBeTruthy();
  });

  it("チップ押下時にonSelectが正しい値で呼ばれる", () => {
    const onSelect = jest.fn();
    const { getByText } = render(
      <ChipSelector items={items} selected="" onSelect={onSelect} />
    );
    fireEvent.press(getByText("浅煎り"));
    expect(onSelect).toHaveBeenCalledWith("light");
  });

  it("別のチップ押下時にonSelectが呼ばれる", () => {
    const onSelect = jest.fn();
    const { getByText } = render(
      <ChipSelector items={items} selected="light" onSelect={onSelect} />
    );
    fireEvent.press(getByText("深煎り"));
    expect(onSelect).toHaveBeenCalledWith("dark");
  });
});
