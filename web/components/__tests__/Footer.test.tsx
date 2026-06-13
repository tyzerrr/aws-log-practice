import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Footer from "../Footer";

describe("Footer", () => {
  it("renders the brand and tagline", () => {
    render(<Footer />);

    expect(screen.getByText("Dumazon")).not.toBeNull();
    expect(
      screen.getByText("暮らしと装いの、ちょうどいい道具を。"),
    ).not.toBeNull();
  });

  it("renders section headings", () => {
    render(<Footer />);

    for (const heading of ["ショップ", "サポート", "会社情報"]) {
      expect(screen.getByText(heading)).not.toBeNull();
    }
  });

  it("renders footer navigation items", () => {
    render(<Footer />);

    for (const item of [
      "新着商品",
      "アパレル",
      "雑貨",
      "ガジェット",
      "配送について",
      "返品・交換",
      "よくある質問",
      "お問い合わせ",
      "会社概要",
      "プライバシーポリシー",
      "特定商取引法",
    ]) {
      expect(screen.getByText(item)).not.toBeNull();
    }
  });
});
