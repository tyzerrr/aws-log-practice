import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Cart from "../Cart";

describe("Cart", () => {
  it("renders the cart button", () => {
    render(<Cart />);

    const cartButton = screen.getByRole("button", { name: "cart-button" });
    expect(cartButton).not.toBeNull();
    expect(cartButton.getAttribute("type")).toBe("button");
  });
});
