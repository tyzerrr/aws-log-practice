"use client";

import { hashKey } from "@/lib/key";
import { rankingHandlers } from "@/lib/lib";
import type { RankingOrder, RankingSelectorContext } from "@/types";

function rankingSelector(focused: boolean): string {
  let css = `border rounded-full px-3 py-1`;
  if (focused) {
    css = `rounded-full px-3 py-1 text-white bg-black border-none`;
  }
  return css;
}

function getRankingSelectorDisplySet(): Array<{
  display: string;
  order: RankingOrder;
}> {
  return [
    { display: "新しい順", order: "newer" },
    { display: "価格が安い順", order: "lowerPrice" },
    { display: "価格が高い順", order: "higherPrice" },
  ];
}

interface RankingSelectorProps {
  context: RankingSelectorContext;
  onSelectOrder: (selector: RankingSelectorContext) => void;
}

export default function RankingSelector({
  context,
  onSelectOrder,
}: RankingSelectorProps) {
  return (
    <div className="flex justify-between items-center mt-6">
      <div className="flex flex-col gap-1">
        <div className="text-primary-line-strong font-bold text-sm">
          ALL PRODUCTS
        </div>
        <div className="font-bold text-xl font-serif">商品一覧</div>
      </div>
      <div className="flex ites-center gap-6">
        {getRankingSelectorDisplySet().map((value, index) => {
          return (
            <button
              key={hashKey("ranking-selector", value.order)}
              type="button"
              onClick={() =>
                onSelectOrder({
                  index: index,
                  order: value.order,
                  handler: rankingHandlers[value.order],
                })
              }
              className={rankingSelector(context.index === index)}
            >
              {value.display}
            </button>
          );
        })}
      </div>
    </div>
  );
}
