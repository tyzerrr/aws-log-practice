import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";

describe("Card", () => {
  it("renders the card structure with slot attributes and custom classes", () => {
    render(
      <Card className="custom-card">
        <CardHeader className="custom-card-header">
          <CardTitle className="custom-card-title">Order summary</CardTitle>
          <CardDescription className="custom-card-description">
            Three items
          </CardDescription>
        </CardHeader>
        <CardContent className="custom-card-content">
          <p>Total: $42</p>
        </CardContent>
      </Card>,
    );

    const card = screen
      .getByText("Order summary")
      .closest("[data-slot='card']");
    const header = screen
      .getByText("Order summary")
      .closest("[data-slot='card-header']");
    const title = screen.getByText("Order summary");
    const description = screen.getByText("Three items");
    const content = screen.getByText("Total: $42").parentElement;

    expect(card).not.toBeNull();
    expect(card?.className).toContain("bg-card");
    expect(card?.className).toContain("custom-card");
    expect(header).not.toBeNull();
    expect(header?.className).toContain("custom-card-header");
    expect(title.getAttribute("data-slot")).toBe("card-title");
    expect(title.className).toContain("font-semibold");
    expect(title.className).toContain("custom-card-title");
    expect(description.getAttribute("data-slot")).toBe("card-description");
    expect(description.className).toContain("text-muted-foreground");
    expect(description.className).toContain("custom-card-description");
    expect(content?.getAttribute("data-slot")).toBe("card-content");
    expect(content?.className).toContain("custom-card-content");
  });
});
