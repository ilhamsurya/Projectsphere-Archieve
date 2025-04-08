CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT name_check CHECK (char_length(name) >= 5 AND char_length(name) <= 50),
    CONSTRAINT username_check CHECK (char_length(username) >= 5 AND char_length(username) <= 15)
);

CREATE TABLE IF NOT EXISTS products(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    name VARCHAR NOT NULL,
    price INT NOT NULL,
    stock INT NOT NULL,
    image_url TEXT NOT NULL,
    condition VARCHAR NOT NULL,
    is_purchasable boolean NOT NULL,
    tags text[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bank_accounts (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    bank_name varchar NOT NULL,
    account_name varchar NOT NULL,
    account_number varchar NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NULL,
    user_id uuid NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT bank_account_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS payments (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    product_id uuid NOT NULL,
    bank_account_id uuid NOT NULL,
    quantity int4 NOT NULL,
    payment_proof_image_url varchar NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT payments_pk PRIMARY KEY (id)
);

