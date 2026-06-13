import Image from "next/image";
import { hashKey } from "@/lib/key";

interface ProductCardProps {
  imageURL: string;
  category: string;
  title: string;
  price: number;
  tags: string[];
}

export default function ProductCard({
  imageURL,
  category,
  title,
  price,
  tags,
}: ProductCardProps) {
  const productImageURL = new URL(imageURL);

  return (
    <div className="flex flex-col items-start">
      <Image
        src={productImageURL.href}
        alt={title}
        width={240}
        height={240}
        className="bg-primary-line-strong"
      />
      <div className="flex flex-col items-start gap-3 w-full">
        <div className="text-primary-line-strong text-sm font-light font-serif">
          {category}
        </div>
        <div className="font-serif">{title}</div>
        <div className="flex justify-between items-center w-full">
          <div className="flex items-center gap-1 font-serif font-bold text-sm">
            ￥{price.toLocaleString()}
            <div className="text-sm font-light text-primary-line-strong">
              JPY
            </div>
          </div>
        </div>
        <div className="flex gap-3 items-start">
          {tags.map((tag, index) => {
            return (
              <div
                key={hashKey("product-tag", title, tag, index)}
                className="flex items-center gap-1 bg-green-200 rounded-full px-2 py-1"
              >
                <div className="rounded-full h-2 w-2 bg-green-700" />
                <div className="text-sm font-bold font-serif">{tag}</div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
