-- Create "products" table
CREATE TABLE "products" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" text NOT NULL,
  "description" text NOT NULL DEFAULT '',
  "price_amount" integer NOT NULL,
  "currency_code" character(3) NOT NULL DEFAULT 'JPY',
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "products_price_amount_check" CHECK (price_amount >= 0)
);
-- Create index "idx_products_is_active" to table: "products"
CREATE INDEX "idx_products_is_active" ON "products" ("is_active");
-- Create "transactions" table
CREATE TABLE "transactions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "total_price_amount" integer NOT NULL DEFAULT 0,
  "currency_code" character(3) NOT NULL DEFAULT 'JPY',
  "status" text NOT NULL DEFAULT 'pending',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "transactions_status_check" CHECK (status = ANY (ARRAY['pending'::text, 'paid'::text, 'cancelled'::text, 'failed'::text, 'fulfilled'::text])),
  CONSTRAINT "transactions_total_price_amount_check" CHECK (total_price_amount >= 0)
);
-- Create index "idx_transactions_status_created_at" to table: "transactions"
CREATE INDEX "idx_transactions_status_created_at" ON "transactions" ("status", "created_at" DESC);
-- Create "orders" table
CREATE TABLE "orders" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "transaction_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  "quantity" integer NOT NULL,
  "unit_price_amount" integer NOT NULL,
  "total_price_amount" integer NOT NULL,
  "currency_code" character(3) NOT NULL DEFAULT 'JPY',
  "status" text NOT NULL DEFAULT 'pending',
  "ordered_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "orders_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "orders_transaction_id_fkey" FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "orders_quantity_check" CHECK (quantity > 0),
  CONSTRAINT "orders_status_check" CHECK (status = ANY (ARRAY['pending'::text, 'paid'::text, 'cancelled'::text, 'failed'::text, 'fulfilled'::text])),
  CONSTRAINT "orders_total_price_amount_check" CHECK (total_price_amount >= 0),
  CONSTRAINT "orders_unit_price_amount_check" CHECK (unit_price_amount >= 0)
);
-- Create index "idx_orders_product_id" to table: "orders"
CREATE INDEX "idx_orders_product_id" ON "orders" ("product_id");
-- Create index "idx_orders_status" to table: "orders"
CREATE INDEX "idx_orders_status" ON "orders" ("status");
-- Create index "idx_orders_transaction_id" to table: "orders"
CREATE INDEX "idx_orders_transaction_id" ON "orders" ("transaction_id");
-- Create "stocks" table
CREATE TABLE "stocks" (
  "product_id" uuid NOT NULL,
  "quantity" integer NOT NULL DEFAULT 0,
  "reserved_quantity" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("product_id"),
  CONSTRAINT "stocks_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "stocks_check" CHECK (reserved_quantity <= quantity),
  CONSTRAINT "stocks_quantity_check" CHECK (quantity >= 0),
  CONSTRAINT "stocks_reserved_quantity_check" CHECK (reserved_quantity >= 0)
);
