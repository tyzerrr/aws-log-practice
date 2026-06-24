import type { Product } from "@/types";

export function chunk<T>(items: T[], chunkSize: number): T[][] {
  const results: T[][] = [];

  for (let i = 0; i < items.length; i += chunkSize) {
    results.push(items.slice(i, i + chunkSize));
  }

  return results;
}

export const rankingHandlers = {
  newer: (a: Product, b: Product) => {
    return b.registeredAt.getTime() - a.registeredAt.getTime();
  },

  higherPrice: (a: Product, b: Product) => {
    return b.price - a.price;
  },

  lowerPrice: (a: Product, b: Product) => {
    return a.price - b.price;
  },
};

// 1. Create Promise that calls resolve() afeter ms
// 2. Promise is returned, status is pending
// 3. Await pending promise until promise would be resolved, that is after ms
// 4. Promise is fulfilled
export function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
