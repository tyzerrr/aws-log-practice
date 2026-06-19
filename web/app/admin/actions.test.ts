import { beforeEach, describe, expect, it, vi } from "vitest";
import { registerProduct } from "./actions";

const registerProductMock = vi.fn();

vi.mock("@/lib/client", () => ({
  productClient: {
    registerProduct: (...args: unknown[]) => registerProductMock(...args),
  },
}));

function buildFormData(
  overrides: Partial<Record<string, string>> = {},
): FormData {
  const defaults: Record<string, string> = {
    name: "Fresh Apple",
    description: "A crisp apple",
    priceAmount: "180",
    currencyCode: "JPY",
    isActive: "on",
  };
  const merged = { ...defaults, ...overrides };
  const formData = new FormData();
  for (const [key, value] of Object.entries(merged)) {
    if (value !== undefined) {
      formData.set(key, value);
    }
  }
  return formData;
}

describe("registerProduct", () => {
  beforeEach(() => {
    registerProductMock.mockReset();
    registerProductMock.mockResolvedValue({
      ok: true,
      payload: { product: { id: "p1" } },
    });
  });

  it("maps form data to a ProductInput and returns the client result", async () => {
    const result = await registerProduct(buildFormData());

    expect(registerProductMock).toHaveBeenCalledWith({
      product: {
        name: "Fresh Apple",
        description: "A crisp apple",
        priceAmount: BigInt(180),
        currencyCode: "JPY",
        isActive: true,
      },
    });
    expect(result).toEqual({ ok: true, payload: { product: { id: "p1" } } });
  });

  it("treats a missing isActive checkbox as false", async () => {
    const formData = buildFormData();
    formData.delete("isActive");

    await registerProduct(formData);

    expect(registerProductMock).toHaveBeenCalledWith(
      expect.objectContaining({
        product: expect.objectContaining({ isActive: false }),
      }),
    );
  });
});
