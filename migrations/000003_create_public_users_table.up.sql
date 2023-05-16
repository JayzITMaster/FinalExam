CREATE TABLE IF NOT EXISTS public_user (
    id serial ,
    users_name text, 
    email citext UNIQUE,
    pu_password_hash bytea ,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);