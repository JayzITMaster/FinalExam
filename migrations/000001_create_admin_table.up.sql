CREATE TABLE IF NOT EXISTS admin_users(
    id serial ,
    users_name text, 
    email citext UNIQUE,
    au_password_hash bytea ,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);