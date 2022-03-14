insert into users(email, first_name, last_name, password, is_admin)
values ('t1@example.com', 't1', 'T1',
        '$argon2i$v=19$m=32768,t=4,p=4$NmO0k7mMgSoN13ajp0xM9A$6JAuu2yYz8WLEQCrsG4QY4KwEv13iIvxKle4XcNNAoQ', false),
       ('t2@example.com', 't2', 'T2',
        '$argon2i$v=19$m=32768,t=4,p=4$NmO0k7mMgSoN13ajp0xM9A$6JAuu2yYz8WLEQCrsG4QY4KwEv13iIvxKle4XcNNAoQ', true),
       ('t3@example.com', 't3', 'T3',
        '$argon2i$v=19$m=32768,t=4,p=4$NmO0k7mMgSoN13ajp0xM9A$6JAuu2yYz8WLEQCrsG4QY4KwEv13iIvxKle4XcNNAoQ', false),
       ('t4@example.com', 't4', 'T4',
        '$argon2i$v=19$m=32768,t=4,p=4$NmO0k7mMgSoN13ajp0xM9A$6JAuu2yYz8WLEQCrsG4QY4KwEv13iIvxKle4XcNNAoQ', false)
;

insert into role_types(name)
values ('public'),
       ('private');

insert into roles(name, type)
values ('member', 1),
       ('manager', 1);

insert into permissions(name)
values ('roles_get_self'),
       ('roles_get'),
       ('redundant_1'),
       ('permissions_get');

insert into roles_permissions(role_id, permission_id)
values (1, 1),
       (2, 2)
;

insert into users_roles(user_id, role_id)
values (1, 1),
       (3, 2),
       (4, 1),
       (4, 2)
;