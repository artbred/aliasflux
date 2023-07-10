CREATE TYPE platform AS ENUM (
    'domain'
);

CREATE TYPE sender AS ENUM (
    'user',
    'ai',
    'ai_system',
    'system'
);

create table chats (
    id uuid unique primary key not null,
    user_id uuid not null references users(id),

    platform platform not null,
    settings jsonb not null,

    feature_user_messages smallint not null default 0,

    created_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone null
);

create table chats_messages(
    id bigserial unique primary key not null,
    chat_id uuid not null,

    message text not null,
    sender sender not null,

    created_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone null,

    foreign key (chat_id) references chats(id)
);
