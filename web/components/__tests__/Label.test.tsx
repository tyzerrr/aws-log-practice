import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Input } from "../ui/input";
import { Label } from "../ui/label";

describe("Label", () => {
  it("labels its associated control and forwards class names", () => {
    render(
      <>
        <Label className="custom-label" htmlFor="display-name">
          Display name
        </Label>
        <Input id="display-name" />
      </>,
    );

    const input = screen.getByLabelText("Display name");
    const label = screen.getByText("Display name");
    expect(input.getAttribute("id")).toBe("display-name");
    expect(label.tagName).toBe("LABEL");
    expect(label.getAttribute("data-slot")).toBe("label");
    expect(label.getAttribute("for")).toBe("display-name");
    expect(label.className).toContain("text-foreground");
    expect(label.className).toContain("custom-label");
  });
});
