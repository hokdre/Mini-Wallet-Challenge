# Mini-Wallet-Challenge

Linter:

```
golangci-lint run -c .golangci.yaml
```

## How to run : 

1. please create postgre database and table in the local computer

   ```

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
   ```
2. export some env :

   ```
   REST_HOST=localhost
   REST_PORT=9001
   REST_WRITE_TIMEOUT_IN_SECOND=2m
   REST_READ_TIMEOUT_IN_SECOND=2m

   POSTGRE_HOST=localhost
   POSTGRE_PORT=5432
   POSTGRE_USERNAME=postgres
   POSTGRE_PASSWORD=abcd
   POSTGRE_DB=wallet
   POSTGRE_SSL_MODE=disable
   POSTGRE_MAX_IDLE_CONN=5
   POSTGRE_MAX_OPEN_CONN=40

   AES_SECRET=1111222233334444 # make sure secret 16 character
   ```
3. running :

   ```
   go run ./cmd/rest/main.go
   ```
