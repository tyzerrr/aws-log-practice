CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price_amount INTEGER NOT NULL CHECK (price_amount >= 0),
    currency_code CHAR(3) NOT NULL DEFAULT 'JPY',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE stocks (
    product_id UUID PRIMARY KEY REFERENCES products (id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved_quantity INTEGER NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (reserved_quantity <= quantity)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    total_price_amount INTEGER NOT NULL DEFAULT 0 CHECK (total_price_amount >= 0),
    currency_code CHAR(3) NOT NULL DEFAULT 'JPY',
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'paid', 'cancelled', 'failed', 'fulfilled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions (id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products (id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price_amount INTEGER NOT NULL CHECK (unit_price_amount >= 0),
    total_price_amount INTEGER NOT NULL CHECK (total_price_amount >= 0),
    currency_code CHAR(3) NOT NULL DEFAULT 'JPY',
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'paid', 'cancelled', 'failed', 'fulfilled')),
    ordered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_products_is_active ON products (is_active);
CREATE INDEX idx_transactions_status_created_at ON transactions (status, created_at DESC);
CREATE INDEX idx_orders_transaction_id ON orders (transaction_id);
CREATE INDEX idx_orders_product_id ON orders (product_id);
CREATE INDEX idx_orders_status ON orders (status);
