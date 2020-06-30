CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    description VARCHAR(1000)
);

CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    pid INTEGER REFERENCES projects ON DELETE CASCADE,
    seqNo INTEGER NOT NULL DEFAULT 0,
    name VARCHAR(255) NOT NULL,
    UNIQUE(pid, seqNo),
    UNIQUE(pid, name)
);

CREATE TABLE priorities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(22) NOT NULL
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    status_id INTEGER REFERENCES statuses ON DELETE RESTRICT,
    project_id INTEGER REFERENCES projects ON DELETE CASCADE,
    priority_id INTEGER REFERENCES priorities ON DELETE SET NULL,
    name VARCHAR(500) NOT NULL,
    description VARCHAR(5000)
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES tasks,
    text VARCHAR(5000) NOT NULL,
    modified CURRENT_TIMESTAMP
);
