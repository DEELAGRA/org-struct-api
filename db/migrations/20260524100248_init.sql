-- +goose Up
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) CHECK (char_length(name) BETWEEN 1 AND 200),
    parent_id INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_parent
        FOREIGN KEY (parent_id)
        REFERENCES departments (id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE UNIQUE INDEX idx_unique_name_parent_not_null
    ON departments (name, parent_id)
    WHERE   parent_id IS NOT NULL;

CREATE UNIQUE INDEX idx_unique_name_parent_null
    ON departments (name)
    WHERE parent_id IS NULL;

CREATE Table employees (
    id SERIAL PRIMARY KEY,
    departament_id INTEGER NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT  NULL,
    hired_at DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_departament
        FOREIGN KEY (departament_id)
        REFERENCES departments (id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);




-- +goose Down
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;
