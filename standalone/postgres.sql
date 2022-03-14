BEGIN;
create table users
(
    id         serial primary key,
    email      text      not null unique,
    first_name text,
    last_name  text,
    password   text,
    is_active  boolean   not null default true,
    is_admin   boolean   not null default false,
    created_at timestamp not null default now()
);

create table role_types
(
    id   serial primary key,
    name text not null unique
);

create table roles
(
    id   serial primary key,
    name text not null,
    type integer
        constraint fk_roles_role_types references role_types,
    constraint uq_roles unique (name, type)
);

create table permissions
(
    id   serial primary key,
    name text not null unique
);

create table users_roles
(
    id      serial primary key,
    user_id integer
        constraint fk_users_roles_users
            references users,
    role_id integer
        constraint fk_users_roles_roles
            references roles,
    constraint uq_users_roles unique (user_id, role_id)
);

create table roles_permissions
(
    id            serial primary key,
    role_id       integer
        constraint fk_role_permissions_roles references roles,
    permission_id integer
        constraint fk_role_permissions_permissions references permissions
);
COMMIT;