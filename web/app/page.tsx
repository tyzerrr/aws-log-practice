import ProductTable from "@/components/ProductTable";
import RankingSelector from "@/components/RankingSelector";

export interface Product {
  imageURL: URL;
  category: string;
  title: string;
  price: number;
  tags: string[];
}

const products: Product[] = [
  {
    imageURL: new URL(
      "https://cdn.dummyjson.com/product-images/beauty/essence-mascara-lash-princess/thumbnail.webp",
    ),
    category: "Beauty",
    title: "Essence Mascara Lash Princess",
    price: 1490,
    tags: ["beauty", "mascara"],
  },
  {
    imageURL: new URL(
      "https://cdn.dummyjson.com/product-images/fragrances/calvin-klein-ck-one/thumbnail.webp",
    ),
    category: "Fragrances",
    title: "Calvin Klein CK One",
    price: 7400,
    tags: ["fragrance", "unisex"],
  },
  {
    imageURL: new URL(
      "https://cdn.dummyjson.com/product-images/furniture/annibale-colombo-bed/thumbnail.webp",
    ),
    category: "Furniture",
    title: "Annibale Colombo Bed",
    price: 278000,
    tags: ["furniture", "bed"],
  },
  {
    imageURL: new URL(
      "https://cdn.dummyjson.com/product-images/groceries/apple/thumbnail.webp",
    ),
    category: "Groceries",
    title: "Fresh Apple",
    price: 180,
    tags: ["food", "fruit"],
  },
  {
    imageURL: new URL(
      "https://cdn.dummyjson.com/product-images/home-decoration/decoration-swing/thumbnail.webp",
    ),
    category: "Home Decoration",
    title: "Decoration Swing",
    price: 8990,
    tags: ["interior", "swing"],
  },
  {
    imageURL: new URL(
      "https://cdn.dummyjson.com/product-images/kitchen-accessories/black-whisk/thumbnail.webp",
    ),
    category: "Kitchen Accessories",
    title: "Black Whisk",
    price: 1290,
    tags: ["kitchen", "tool"],
  },
  {
    imageURL: new URL(
      "https://cdn.dummyjson.com/product-images/laptops/apple-macbook-pro-14-inch-space-grey/thumbnail.webp",
    ),
    category: "Laptops",
    title: "Apple MacBook Pro 14 Inch",
    price: 298000,
    tags: ["laptop", "apple"],
  },
  {
    imageURL: new URL(
      "https://cdn.dummyjson.com/product-images/mens-shirts/blue-&-black-check-shirt/thumbnail.webp",
    ),
    category: "Mens Shirts",
    title: "Blue & Black Check Shirt",
    price: 3490,
    tags: ["fashion", "shirt"],
  },
];

export default function Home() {
  return (
    <div className="w-full flex flex-col justify-between px-12">
      <RankingSelector />
      <ProductTable
        // TODO: Need to calculate dynamically
        tableLayout={{
          totalProductCount: products.length,
          column: 4,
          row: Math.ceil(products.length / 4),
        }}
        products={products}
      />
    </div>
  );
}
