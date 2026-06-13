import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { rankingHandlers } from "@/lib/lib";
import type { RankingSelectorContext } from "@/types";
import RankingSelector from "../RankingSelector";

describe("RankingSelector", () => {
  const context: RankingSelectorContext = {
    index: 0,
    order: "newer",
    handler: rankingHandlers.newer,
  };

  it("renders the product listing heading and order buttons", () => {
    render(<RankingSelector context={context} onSelectOrder={vi.fn()} />);

    expect(screen.getByText("ALL PRODUCTS")).not.toBeNull();
    expect(screen.getByText("商品一覧")).not.toBeNull();

    for (const label of ["新しい順", "価格が安い順", "価格が高い順"]) {
      expect(screen.getByRole("button", { name: label })).not.toBeNull();
    }
  });

  it("marks the focused selector from context", () => {
    render(
      <RankingSelector
        context={{
          index: 1,
          order: "lowerPrice",
          handler: rankingHandlers.lowerPrice,
        }}
        onSelectOrder={vi.fn()}
      />,
    );

    expect(
      screen
        .getByRole("button", { name: "価格が安い順" })
        .className.includes("bg-black"),
    ).toBe(true);
    expect(
      screen
        .getByRole("button", { name: "新しい順" })
        .className.includes("bg-black"),
    ).toBe(false);
  });

  it("calls onSelectOrder with the selected ranking context", () => {
    const onSelectOrder = vi.fn();
    render(<RankingSelector context={context} onSelectOrder={onSelectOrder} />);

    fireEvent.click(screen.getByRole("button", { name: "価格が高い順" }));

    expect(onSelectOrder).toHaveBeenCalledOnce();
    expect(onSelectOrder).toHaveBeenCalledWith({
      index: 2,
      order: "higherPrice",
      handler: rankingHandlers.higherPrice,
    });
  });
});
