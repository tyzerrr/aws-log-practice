import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Textarea } from "../ui/textarea";

describe("Textarea", () => {
  it("renders a textarea with forwarded attributes and class names", () => {
    render(
      <Textarea
        aria-label="Message"
        className="custom-textarea"
        defaultValue="Initial message"
        placeholder="Write a message"
        rows={4}
      />,
    );

    const textarea = screen.getByRole("textbox", {
      name: "Message",
    }) as HTMLTextAreaElement;
    expect(textarea.tagName).toBe("TEXTAREA");
    expect(textarea.getAttribute("data-slot")).toBe("textarea");
    expect(textarea.value).toBe("Initial message");
    expect(textarea.getAttribute("placeholder")).toBe("Write a message");
    expect(textarea.getAttribute("rows")).toBe("4");
    expect(textarea.className).toContain("min-h-24");
    expect(textarea.className).toContain("custom-textarea");
  });
});
