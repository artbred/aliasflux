DO
$do$
    declare
        _key text;
        _value jsonb;
    begin
        create table if not exists settings (
            key text unique not null,
            value jsonb not null
        );

        _key := 'platform';
        _value := '[
            {"platform": "domain"}
        ]'::jsonb;

        if not exists (select 1 from settings where key = _key) then
            insert into settings (key, value) values (_key, _value);
        else
            update settings set value = _value where key = _key;
        end if;
    end
$do$;



