"use server";

import { productClient } from "@/lib/client";

export async function registerProduct(formData: FormData) {
  // TODO(auth): verify the caller is an authenticated admin once auth lands.

  // The retryable client handles retries and Result wrapping, so the action
  // only maps the form data to a ProductInput and delegates the RPC call.
  return productClient.registerProduct({
    product: {
      name: String(formData.get("name") ?? ""),
      description: String(formData.get("description") ?? ""),
      priceAmount: BigInt(String(formData.get("priceAmount") ?? "0") || "0"),
      currencyCode: String(formData.get("currencyCode") ?? ""),
      // An unchecked checkbox is omitted from the form data entirely.
      isActive: formData.get("isActive") === "on",
    },
  });
}
