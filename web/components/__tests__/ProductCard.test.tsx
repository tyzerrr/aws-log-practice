import { render, screen } from "@testing-library/react";
import type { ImgHTMLAttributes } from "react";
import { describe, expect, it, vi } from "vitest";
import ProductCard from "../ProductCard";

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

describe("ProductCard", () => {
  it("renders product image, category, title, and price", () => {
    render(
      <ProductCard
        imageURL="https://example.com/products/mug.webp"
        category="Kitchen"
        title="Ceramic Mug"
        price={12800}
        tags={["drinkware", "gift"]}
      />,
    );

    const image = screen.getByRole("img", { name: "Ceramic Mug" });
    expect(image.getAttribute("src")).toBe(
      "https://example.com/products/mug.webp",
    );
    expect(screen.getByText("Kitchen")).not.toBeNull();
    expect(screen.getByText("Ceramic Mug")).not.toBeNull();
    const currency = screen.getByText("JPY");
    expect(currency).not.toBeNull();
    expect(currency.parentElement?.textContent).toBe("￥12,800JPY");
  });

  it("renders all product tags", () => {
    render(
      <ProductCard
        imageURL="https://example.com/products/jacket.webp"
        category="Fashion"
        title="Rain Jacket"
        price={9400}
        tags={["waterproof", "outerwear", "travel"]}
      />,
    );

    for (const tag of ["waterproof", "outerwear", "travel"]) {
      expect(screen.getByText(tag)).not.toBeNull();
    }
  });
});
