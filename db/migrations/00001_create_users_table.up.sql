CREATE TABLE IF NOT EXISTS auth.users (
                                     id SERIAL PRIMARY KEY,
                                     email VARCHAR(255) NOT NULL,
                                     name VARCHAR(255) NOT NULL,
                                     password BYTEA,
                                     created_at TIMESTAMP DEFAULT NOW(),
                                     updated_at TIMESTAMP DEFAULT NOW(),
                                     deleted_at TIMESTAMP DEFAULT NOW(),
                                     active BOOL NOT NULL
);