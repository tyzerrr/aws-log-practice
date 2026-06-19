"use client";

import ProductTable from "@/components/ProductTable";
import { useSearch } from "@/components/SearchProvider";
import type { Product } from "@/types";
import { useState } from "react";

const products: Product[] = [
  {
    imageURL:
      "https://cdn.dummyjson.com/product-images/beauty/essence-mascara-lash-princess/thumbnail.webp",
    category: "Beauty",
    title: "Essence Mascara Lash Princess",
    price: 1490,
    tags: ["beauty", "mascara"],
    registeredAt: new Date("2026-06-12T01:15:00.000Z"),
  },
  {
    imageURL:
      "https://cdn.dummyjson.com/product-images/fragrances/calvin-klein-ck-one/thumbnail.webp",
    category: "Fragrances",
    title: "Calvin Klein CK One",
    price: 7400,
    tags: ["fragrance", "unisex"],
    registeredAt: new Date("2026-06-11T08:40:00.000Z"),
  },
  {
    imageURL:
      "https://cdn.dummyjson.com/product-images/furniture/annibale-colombo-bed/thumbnail.webp",
    category: "Furniture",
    title: "Annibale Colombo Bed",
    price: 278000,
    tags: ["furniture", "bed"],
    registeredAt: new Date("2026-06-10T13:05:00.000Z"),
  },
  {
    imageURL:
      "https://cdn.dummyjson.com/product-images/groceries/apple/thumbnail.webp",
    category: "Groceries",
    title: "Fresh Apple",
    price: 180,
    tags: ["food", "fruit"],
    registeredAt: new Date("2026-06-09T22:20:00.000Z"),
  },
  {
    imageURL:
      "https://cdn.dummyjson.com/product-images/home-decoration/decoration-swing/thumbnail.webp",
    category: "Home Decoration",
    title: "Decoration Swing",
    price: 8990,
    tags: ["interior", "swing"],
    registeredAt: new Date("2026-06-08T04:55:00.000Z"),
  },
  {
    imageURL:
      "https://cdn.dummyjson.com/product-images/kitchen-accessories/black-whisk/thumbnail.webp",
    category: "Kitchen Accessories",
    title: "Black Whisk",
    price: 1290,
    tags: ["kitchen", "tool"],
    registeredAt: new Date("2026-06-07T17:30:00.000Z"),
  },
  {
    imageURL:
      "https://cdn.dummyjson.com/product-images/laptops/apple-macbook-pro-14-inch-space-grey/thumbnail.webp",
    category: "Laptops",
    title: "Apple MacBook Pro 14 Inch",
    price: 298000,
    tags: ["laptop", "apple"],
    registeredAt: new Date("2026-06-06T10:10:00.000Z"),
  },
  {
    imageURL:
      "https://cdn.dummyjson.com/product-images/mens-shirts/blue-&-black-check-shirt/thumbnail.webp",
    category: "Mens Shirts",
    title: "Blue & Black Check Shirt",
    price: 3490,
    tags: ["fashion", "shirt"],
    registeredAt: new Date("2026-06-05T15:45:00.000Z"),
  },
];

export default function Home() {
  const { searchTerm } = useSearch();
  return (
    <div className="w-full flex flex-col justify-between px-12">
      <ProductTable
        // TODO: Need to calculate dynamically
        tableLayout={{
          totalProductCount: products.length,
          column: 4,
          row: Math.ceil(products.length / 4),
        }}
        products={products}
        searchTerm={searchTerm}
      />
    </div>
  );
}
