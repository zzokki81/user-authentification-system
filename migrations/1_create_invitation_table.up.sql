CREATE TABLE invitations(
    id serial primary key,
    email varchar(255) not null unique,
    inviter_id int,
    created_at timestamptz not null
);