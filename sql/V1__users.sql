create table users(
    id uuid unique primary key,
    create_params jsonb null,
    created_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone null
)
