import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Badge } from "../ui/badge";

describe("Badge", () => {
  it("renders a default badge with forwarded props", () => {
    render(<Badge className="custom-badge">New</Badge>);

    const badge = screen.getByText("New");
    expect(badge.tagName).toBe("SPAN");
    expect(badge.getAttribute("data-slot")).toBe("badge");
    expect(badge.className).toContain("bg-primary");
    expect(badge.className).toContain("custom-badge");
  });

  it("applies variant classes", () => {
    render(<Badge variant="success">Active</Badge>);

    const badge = screen.getByText("Active");
    expect(badge.className).toContain("border-emerald-200");
    expect(badge.className).toContain("text-emerald-700");
  });
});
