import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Header from "../Header";
import Likes from "../Likes";
import Search from "../Search";

describe("Header", () => {
  it("renders the logo link to home", () => {
    render(<Header />);
    const logo = screen.getByRole("link", { name: "Dumazon" });
    expect(logo).not.toBeNull();
    expect(logo.getAttribute("href")).toBe("/");
  });
});

describe("Header", () => {
  it("renders the likes to home", () => {
    render(<Likes />);
    const likes = screen.getByRole("button", { name: "likes-button" });
    expect(likes).not.toBeNull();
  });
});

describe("Header", () => {
  it("renders the search to home", () => {
    render(<Search />);
    const searchInput = screen.getByRole("textbox", { name: "search-input" });
    expect(searchInput).not.toBeNull();

    const placeholder = screen.getByPlaceholderText("商品を検索...");
    expect(placeholder).not.toBeNull();
  });
});
