DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    role VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    firstname VARCHAR(100) NOT NULL
);

DROP TABLE IF EXISTS groups CASCADE;
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100)
);

DROP TABLE IF EXISTS users_in_groups;
CREATE TABLE users_in_groups (
    id SERIAL PRIMARY KEY,
    userid INTEGER NOT NULL,
    groupid INTEGER NOT NULL,
    FOREIGN KEY (userid) REFERENCES users(id),
    FOREIGN KEY (groupid) REFERENCES groups(id)
);

DROP TABLE IF EXISTS contests CASCADE;
CREATE TABLE contests (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    time TIMESTAMP,
    duration INTEGER,
    groupid INTEGER NOT NULL,
    FOREIGN KEY (groupid) REFERENCES groups(id)
);

DROP TABLE IF EXISTS tasks CASCADE;
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    contestid INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    url VARCHAR(255) NOT NULL,
    FOREIGN KEY (contestid) REFERENCES contests(id)
);

DROP TABLE IF EXISTS submissions;
CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    userid INTEGER NOT NULL,
    taskid INTEGER NOT NULL,
    status VARCHAR(100) NOT NULL,
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    FOREIGN KEY (userid) REFERENCES users(id),
    FOREIGN KEY (taskid) REFERENCES tasks(id)
);
