CREATE TABLE IF NOT EXISTS users (
 user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
 email VARCHAR(150) NOT NULL UNIQUE,
 password_hash BYTEA NOT NULL,
 is_superuser BOOLEAN DEFAULT FALSE,
 first_name VARCHAR(255) NOT NULL,
 last_name VARCHAR(255) NOT NULL
);