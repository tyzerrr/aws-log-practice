import Image from "next/image";

interface ProductCardProps {
  imageURL: URL;
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
  return (
    <div className="flex flex-col items-center">
      <Image src={imageURL.href} alt={`${title}`} width={100} height={300} />
      <div className="flex flex-col items-center gap-3">
        <div className="text-primary-line-strong text-sm font-light font-serif">
          {category}
        </div>
        <div className="font-serif">{title}</div>
        <div className="flex justify-between items-center">
          <div className="font-serif font-bold">
            `￥${price}`
            <span className="text-sm font-light text-primary-line-strong">
              JPY
            </span>
          </div>
          <div className="flex gap-3 items-center bg-green-400 text-green-900">
            {tags.map((tag) => {
              return <span className="">{tag}</span>;
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
