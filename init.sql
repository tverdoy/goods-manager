CREATE TABLE IF NOT EXISTS projects
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

INSERT INTO projects (name)
VALUES ('Project 1');

CREATE TABLE IF NOT EXISTS goods
(
    id          SERIAL PRIMARY KEY,
    project_id  INT REFERENCES projects (id),
    name        VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    priority    INT          NOT NULL,
    removed     BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_goods_name ON goods (name);