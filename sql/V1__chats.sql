CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SET TIME ZONE 'UTC';

create table chats (
    id uuid unique primary key not null,
    created_at timestamp with time zone not null,
    deleted_at timestamp with time zone null
);

create table chats_messages(
    id bigserial unique primary key not null,
    chat_id uuid not null,
    message text not null,
    is_system boolean not null default false,
    created_at timestamp with time zone not null,
    deleted_at timestamp with time zone null,
    foreign key (chat_id) references chats(id)
);

create table chats_configs(
    id bigserial unique primary key not null,
    chat_id uuid not null,
    config jsonb not null,
    created_at timestamp with time zone not null,
    foreign key (chat_id) references chats(id)
);
