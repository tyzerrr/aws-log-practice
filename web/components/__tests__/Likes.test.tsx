import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Likes from "../Likes";

describe("Likes", () => {
  it("renders the likes button", () => {
    render(<Likes />);

    const likesButton = screen.getByRole("button", { name: "likes-button" });
    expect(likesButton).not.toBeNull();
    expect(likesButton.getAttribute("type")).toBe("button");
  });
});
