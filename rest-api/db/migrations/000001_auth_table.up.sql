
create table api_users (
    id serial primary key,
    name varchar(40) UNIQUE,
    username varchar(40),
    api_key varchar(40)
);
