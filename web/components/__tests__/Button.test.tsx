import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { Button } from "../ui/button";

describe("Button", () => {
  it("renders a button with default classes and forwarded props", () => {
    const handleClick = vi.fn();

    render(
      <Button
        aria-label="Save"
        className="custom-button"
        type="button"
        onClick={handleClick}
      >
        Save
      </Button>,
    );

    const button = screen.getByRole("button", { name: "Save" });
    fireEvent.click(button);

    expect(button.getAttribute("data-slot")).toBe("button");
    expect(button.getAttribute("type")).toBe("button");
    expect(button.className).toContain("bg-primary");
    expect(button.className).toContain("h-9");
    expect(button.className).toContain("custom-button");
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it("applies variant and size classes", () => {
    render(
      <Button variant="outline" size="sm">
        Filter
      </Button>,
    );

    const button = screen.getByRole("button", { name: "Filter" });
    expect(button.className).toContain("border-input");
    expect(button.className).toContain("h-8");
    expect(button.className).toContain("px-3");
  });

  it("renders its child element when asChild is set", () => {
    render(
      <Button asChild variant="ghost" className="custom-link-button">
        <a href="/products">Products</a>
      </Button>,
    );

    const link = screen.getByRole("link", { name: "Products" });
    expect(link.tagName).toBe("A");
    expect(link.getAttribute("href")).toBe("/products");
    expect(link.getAttribute("data-slot")).toBe("button");
    expect(link.className).toContain("hover:bg-accent");
    expect(link.className).toContain("custom-link-button");
  });
});
