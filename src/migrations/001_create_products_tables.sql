CREATE TABLE categories
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    parent_id  UUID REFERENCES categories (id),
    is_active  BOOLEAN                  DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE sellers
(
    id           UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    company_name VARCHAR(255)        NOT NULL,
    email        VARCHAR(255) UNIQUE NOT NULL,
    phone        VARCHAR(50),
    address      TEXT,
    tax_number   VARCHAR(50) UNIQUE,
    is_active    BOOLEAN                  DEFAULT true,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE products
(
    id          UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name        VARCHAR(255)        NOT NULL,
    description TEXT,
    category_id UUID                NOT NULL REFERENCES categories (id),
    brand       VARCHAR(100),
    sku         VARCHAR(100) UNIQUE NOT NULL,
    status      VARCHAR(20)              DEFAULT 'active',
    images      TEXT[],
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at  TIMESTAMP WITH TIME ZONE
);

CREATE TABLE product_variants
(
    id             UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    product_id     UUID           NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    seller_id      UUID           NOT NULL REFERENCES sellers (id) ON DELETE CASCADE,
    price          DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    discount_price DECIMAL(10, 2) CHECK (discount_price IS NULL OR discount_price > 0),
    stock          INTEGER        NOT NULL  DEFAULT 0 CHECK (stock >= 0),
    attributes     JSONB,
    is_active      BOOLEAN                  DEFAULT true,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_products_category_id ON products (category_id);
CREATE INDEX idx_products_sku ON products (sku);
CREATE INDEX idx_products_status ON products (status);
CREATE INDEX idx_products_name ON products (name);
CREATE INDEX idx_products_deleted_at ON products (deleted_at);

CREATE INDEX idx_product_variants_product_id ON product_variants (product_id);
CREATE INDEX idx_product_variants_seller_id ON product_variants (seller_id);
CREATE INDEX idx_product_variants_price ON product_variants (price);
CREATE INDEX idx_product_variants_stock ON product_variants (stock);
CREATE INDEX idx_product_variants_is_active ON product_variants (is_active);

CREATE INDEX idx_categories_parent_id ON categories (parent_id);
CREATE INDEX idx_categories_is_active ON categories (is_active);