import { fireEvent, render, screen } from "@testing-library/react";
import type { ImgHTMLAttributes } from "react";
import { describe, expect, it, vi } from "vitest";
import type { Product } from "@/types";
import ProductTable from "../ProductTable";

vi.mock("next/image", () => ({
  default: ({
    src,
    alt,
    ...props
  }: ImgHTMLAttributes<HTMLImageElement> & { src: string }) => {
    // biome-ignore lint/performance/noImgElement: next/image is mocked as a native image for DOM assertions.
    return <img src={src} alt={alt} {...props} />;
  },
}));

const products: Product[] = [
  {
    imageURL: "https://example.com/products/mid.webp",
    category: "Home",
    title: "Middle Product",
    price: 3000,
    tags: ["middle"],
    registeredAt: new Date("2026-06-10T00:00:00.000Z"),
  },
  {
    imageURL: "https://example.com/products/old.webp",
    category: "Home",
    title: "Oldest Product",
    price: 1000,
    tags: ["old"],
    registeredAt: new Date("2026-06-01T00:00:00.000Z"),
  },
  {
    imageURL: "https://example.com/products/new.webp",
    category: "Home",
    title: "Newest Product",
    price: 5000,
    tags: ["new"],
    registeredAt: new Date("2026-06-12T00:00:00.000Z"),
  },
];

function renderProductTable() {
  return render(
    <ProductTable
      tableLayout={{
        row: 2,
        column: 2,
        totalProductCount: products.length,
      }}
      products={products}
      searchTerm=""
    />,
  );
}

function renderedProductTitles() {
  return screen.getAllByRole("img").map((image) => image.getAttribute("alt"));
}

describe("ProductTable", () => {
  it("renders selector labels and products sorted by newest first", () => {
    renderProductTable();

    expect(screen.getByText("ALL PRODUCTS")).not.toBeNull();
    expect(screen.getByText("商品一覧")).not.toBeNull();
    expect(renderedProductTitles()).toEqual([
      "Newest Product",
      "Middle Product",
      "Oldest Product",
    ]);
  });

  it("sorts products when a ranking selector is clicked", () => {
    renderProductTable();

    fireEvent.click(screen.getByRole("button", { name: "価格が安い順" }));
    expect(renderedProductTitles()).toEqual([
      "Oldest Product",
      "Middle Product",
      "Newest Product",
    ]);

    fireEvent.click(screen.getByRole("button", { name: "価格が高い順" }));
    expect(renderedProductTitles()).toEqual([
      "Newest Product",
      "Middle Product",
      "Oldest Product",
    ]);
  });
});
