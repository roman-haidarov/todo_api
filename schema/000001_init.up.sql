CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    username      varchar(255) not null,
    tsv           tsvector GENERATED ALWAYS AS (
                    to_tsvector('english', coalesce(username, '') || ' ') ||
                    to_tsvector('russian', coalesce(username, '') || ' ')
                  ) STORED,
    password_hash varchar(255) not null
);

CREATE TABLE todo_lists
(
    id           serial       not null unique,
    title        varchar(255) not null,
    descriptions varchar(255) not null,
    tsv          tsvector GENERATED ALWAYS AS (
                    to_tsvector('english', coalesce(title, '') || ' ' || coalesce(descriptions, ''))
                 ) STORED,
    PRIMARY KEY (id)
);

CREATE TABLE user_lists
(
    id      serial                                          not null unique,
    user_id int references users(id) on delete cascade      not null,
    list_id int references todo_lists(id) on delete cascade not null
);

CREATE TABLE todo_items
(
    id           serial       not null unique,
    title        varchar(255) not null,
    descriptions varchar(255),
    done         boolean      not null default false,
    tsv          tsvector GENERATED ALWAYS AS (
                    to_tsvector('english', coalesce(title, '') || ' ' || coalesce(descriptions, '')) ||
                    to_tsvector('russian', coalesce(title, '') || ' ' || coalesce(descriptions, ''))
                 ) STORED,
    PRIMARY KEY (id)
);

CREATE TABLE lists_items
(
    id      serial                                          not null unique,
    item_id int references todo_items(id) on delete cascade not null,
    list_id int references todo_lists(id) on delete cascade not null
);

CREATE TABLE user_items
(
    id      serial                                          not null unique,
    item_id int references todo_items(id) on delete cascade not null,
    user_id int references users(id)      on delete cascade not null
);

CREATE INDEX IF NOT EXISTS idx_todo_lists_tsv ON todo_lists USING gin(tsv);
CREATE INDEX IF NOT EXISTS idx_users_tsv ON users USING gin(tsv);
CREATE INDEX IF NOT EXISTS idx_todo_items_tsv ON todo_items USING gin(tsv);
