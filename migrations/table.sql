
CREATE DATABASE mywallet;


CREATE TABLE accounts (
    id VARCHAR(36) NOT NULL,
    external_id VARCHAR(36) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT 'true' 
    PRIMARY KEY(id)
)


CREATE TABLE wallets (
    id VARCHAR(36) NOT NULL,
    owned_by VARCHAR(36) UNIQUE NOT NULL,
    balance NUMERIC NOT NULL,
    status VARCHAR(255) NOT NULL,
    enabled_at TIMESTAMP NULL,
    disabled_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT 'true',
    PRIMARY KEY(id),
    FOREIGN KEY (owned_by) REFERENCES accounts(id)
)

CREATE TABLE transactions (
    id VARCHAR(36) NOT NULL,
    wallet_id VARCHAR(36) NOT NULL,
    type VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    reference_id VARCHAR(255) UNIQUE NOT NULL,
    amount NUMERIC NOT NULL,
    transacted_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT 'true',
    PRIMARY KEY(id),
    FOREIGN KEY (wallet_id) REFERENCES wallets(id)
)