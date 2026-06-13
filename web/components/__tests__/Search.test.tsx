import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Search from "../Search";

describe("Search", () => {
  it("renders an accessible search input with the expected placeholder", () => {
    render(<Search />);

    const searchInput = screen.getByRole("textbox", { name: "search-input" });
    expect(searchInput).not.toBeNull();
    expect(searchInput.getAttribute("placeholder")).toBe("商品を検索...");
  });
});
