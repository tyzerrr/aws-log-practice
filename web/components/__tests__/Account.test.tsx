import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Account from "../Account";

describe("Account", () => {
  it("renders the account button", () => {
    render(<Account />);

    const accountButton = screen.getByRole("button", {
      name: "account-button",
    });
    expect(accountButton).not.toBeNull();
    expect(accountButton.getAttribute("type")).toBe("button");
  });
});
