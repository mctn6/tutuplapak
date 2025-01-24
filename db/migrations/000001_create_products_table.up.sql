CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    category VARCHAR(32) NOT NULL,
    qty INT NOT NULL CHECK (qty >= 1),
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 100),
    sku VARCHAR(32) NOT NULL,
    fileId INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (fileId) REFERENCES files(id)
);

CREATE INDEX idx_products_name ON products(name);

CREATE INDEX idx_products_fileId ON products (fileId);
CREATE INDEX idx_files_id ON files (id);