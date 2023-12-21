-- +goose Up
-- +goose StatementBegin

create type url_schema as enum ('http', 'https');

create table if not exists link
(
    id         text primary key, -- hash from uuid
    "schema"   url_schema  not null,
    "domain"   text        not null,
    uri        text        not null,
    created_at timestamptz not null default current_timestamp,
    expired_at timestamptz          default null
) partition by hash (id);

create table if not exists link_p0 partition of link for values with (modulus 8, remainder 0);
create table if not exists link_p1 partition of link for values with (modulus 8, remainder 1);
create table if not exists link_p2 partition of link for values with (modulus 8, remainder 2);
create table if not exists link_p3 partition of link for values with (modulus 8, remainder 3);
create table if not exists link_p4 partition of link for values with (modulus 8, remainder 4);
create table if not exists link_p5 partition of link for values with (modulus 8, remainder 5);
create table if not exists link_p6 partition of link for values with (modulus 8, remainder 6);
create table if not exists link_p7 partition of link for values with (modulus 8, remainder 7);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists link;
-- +goose StatementEnd
