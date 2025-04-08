CREATE TABLE IF NOT EXISTS users (
 user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
 email TEXT NOT NULL UNIQUE,
 password_hash BYTEA NOT NULL,
 is_superuser BOOLEAN DEFAULT FALSE,
 first_name TEXT NOT NULL,
 last_name TEXT NOT NULL
);