import { ShoppingCart } from "lucide-react";

export default function Cart() {
  return (
    <button type="button" aria-label="cart-button">
      <ShoppingCart className={"size-6"} />
    </button>
  );
}
