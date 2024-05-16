CREATE TABLE artifact (
    artifact_id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    age INTEGER,
    location VARCHAR(100) NOT NULL,
    researcher_id INTEGER REFERENCES researcher(researcher_id) NOT NULL
);