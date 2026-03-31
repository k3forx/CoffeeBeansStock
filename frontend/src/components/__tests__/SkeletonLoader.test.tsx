import React from "react";
import { render } from "@testing-library/react-native";
import { BeanListSkeleton } from "../SkeletonLoader";

describe("SkeletonLoader", () => {
  it("BeanListSkeletonがクラッシュせずにレンダリングされる", () => {
    expect(() => render(<BeanListSkeleton />)).not.toThrow();
  });
});
