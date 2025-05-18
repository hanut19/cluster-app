DROP TABLE IF EXISTS users, clusters;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('readonly', 'admin'))
);

CREATE TABLE clusters (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    Nodes INTEGER NOT NULL
);

INSERT INTO users (username, password, role)
VALUES ('admin', crypt('admin', gen_salt('bf')), 'admin'),
       ('readonly', crypt('readonly', gen_salt('bf')), 'readonly');

INSERT INTO clusters (name, nodes)
VALUES ('cluster-a', 10), ('cluster-b', 5), ('cluster-c', 3), ('cluster-d', 0);