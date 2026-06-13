import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Input } from "../ui/input";

describe("Input", () => {
  it("renders an input with forwarded attributes and class names", () => {
    render(
      <Input
        aria-label="Email"
        className="custom-input"
        defaultValue="person@example.com"
        placeholder="you@example.com"
        type="email"
      />,
    );

    const input = screen.getByRole("textbox", {
      name: "Email",
    }) as HTMLInputElement;
    expect(input.tagName).toBe("INPUT");
    expect(input.getAttribute("data-slot")).toBe("input");
    expect(input.type).toBe("email");
    expect(input.value).toBe("person@example.com");
    expect(input.getAttribute("placeholder")).toBe("you@example.com");
    expect(input.className).toContain("border-input");
    expect(input.className).toContain("custom-input");
  });
});
