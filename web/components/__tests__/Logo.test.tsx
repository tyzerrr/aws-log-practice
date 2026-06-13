import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Logo from "../Logo";

describe("Logo", () => {
  it("renders a home link with the brand name", () => {
    render(<Logo />);

    const logo = screen.getByRole("link", { name: "Dumazon" });
    expect(logo).not.toBeNull();
    expect(logo.getAttribute("href")).toBe("/");
  });
});
