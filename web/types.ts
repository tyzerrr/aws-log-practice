export interface Product {
  imageURL: string;
  category: string;
  title: string;
  price: number;
  tags: string[];
  registeredAt: Date;
}

export type CompareFn<T> = (a: T, b: T) => number;
export type RankingHandler = CompareFn<Product>;

export type RankingOrder = "newer" | "lowerPrice" | "higherPrice";

export interface RankingSelectorContext {
  index: number;
  order: RankingOrder;
  handler: RankingHandler;
}
