"use client";

import { useState } from "react";
import { hashKey } from "@/lib/key";
import { chunk, rankingHandlers } from "@/lib/lib";
import type { Product, RankingSelectorContext } from "@/types";
import ProductCard from "./ProductCard";
import RankingSelector from "./RankingSelector";

interface ProductTableLayout {
  row: number;
  column: number;
  totalProductCount: number;
}

interface ProductTableProps {
  tableLayout: ProductTableLayout;
  products: Product[];
}

interface ProductTableRowProps {
  products: Product[];
}

function ProductTableRow({ products }: ProductTableRowProps) {
  return (
    <div className="flex gap-6 items-center w-full justify-between">
      {products.map((product, index) => {
        return (
          <ProductCard
            key={hashKey("product", product, index)}
            imageURL={product.imageURL}
            category={product.category}
            title={product.title}
            price={product.price}
            tags={product.tags}
          />
        );
      })}
    </div>
  );
}

export default function ProductTable({
  tableLayout,
  products,
}: ProductTableProps) {
  const [focusedRankingSelectorContext, setFocusedRankingSelectorContext] =
    useState<RankingSelectorContext>({
      index: 0,
      order: "newer",
      handler: rankingHandlers.newer,
    });
  return (
    <div>
      <RankingSelector
        context={focusedRankingSelectorContext}
        onSelectOrder={(selectorContext) =>
          setFocusedRankingSelectorContext(selectorContext)
        }
      />
      {chunk(
        products.toSorted(focusedRankingSelectorContext.handler),
        tableLayout.column,
      ).map((productsChunk, index) => {
        return (
          <div
            key={hashKey("product-row", productsChunk, index)}
            className="flex flex-col justify-between w-full my-6"
          >
            <ProductTableRow products={productsChunk} />
          </div>
        );
      })}
    </div>
  );
}
