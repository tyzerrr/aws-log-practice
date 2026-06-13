import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Separator } from "../ui/separator";

describe("Separator", () => {
  it("renders an hr separator with forwarded props", () => {
    render(<Separator className="custom-separator" data-testid="separator" />);

    const separator = screen.getByTestId("separator");
    expect(separator.tagName).toBe("HR");
    expect(separator.getAttribute("data-slot")).toBe("separator");
    expect(separator.className).toContain("bg-border");
    expect(separator.className).toContain("custom-separator");
  });
});
