CREATE TABLE expedition (
    expedition_id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    expeditionYear INT,
    researcher_id INTEGER REFERENCES researcher(researcher_id) NOT NULL
);