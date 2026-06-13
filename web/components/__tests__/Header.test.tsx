import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Header from "../Header";

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
    render(<Header />);
    const likes = screen.getByRole("button", { name: "likes-button" });
    expect(likes).not.toBeNull();
  });
});

describe("Header", () => {
  it("renders the search to home", () => {
    render(<Header />);
    const searchInput = screen.getByRole("textbox", { name: "search-input" });
    expect(searchInput).not.toBeNull();

    const placeholder = screen.getByPlaceholderText("商品を検索...");
    expect(placeholder).not.toBeNull();
  });
});

describe("Header", () => {
  it("renders the account to home", () => {
    render(<Header />);
    const account = screen.getByRole("button", { name: "account-button" });
    expect(account).not.toBeNull();
  });
});

describe("Header", () => {
  it("renders the cart to home", () => {
    render(<Header />);
    const cart = screen.getByRole("button", { name: "cart-button" });
    expect(cart).not.toBeNull();
  });
});
