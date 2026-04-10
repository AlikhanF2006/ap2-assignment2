CREATE TABLE payments (
                          id TEXT PRIMARY KEY,
                          order_id TEXT NOT NULL,
                          transaction_id TEXT NOT NULL UNIQUE,
                          amount BIGINT NOT NULL CHECK (amount > 0),
                          status TEXT NOT NULL
);