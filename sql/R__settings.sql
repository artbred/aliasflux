DO
$do$
    declare
        _key text;
        _value jsonb;
    begin
        create table if not exists settings (
            platform platform null,
            key text not null primary key,
            value jsonb not null,

            unique(platform, key)
        );

        _key := 'free_chat_features';
        _value := '{
            "messages_limit": 2
        }'::jsonb;

        if not exists (select 1 from settings where key = _key) then
            insert into settings (key, value) values (_key, _value);
        else
            update settings set value = _value where key = _key;
        end if;
    end
$do$;



