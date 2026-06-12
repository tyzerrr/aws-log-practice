interface ProductTableProps {
  productsCount: number;
}

export default function ProductTable({ productsCount }: ProductTableProps) {
  return <div>{productsCount}</div>;
}
