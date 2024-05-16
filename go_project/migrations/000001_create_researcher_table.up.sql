CREATE TABLE researcher (
    researcher_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    specialization VARCHAR(100) NOT NULL,
    project VARCHAR(100) NOT NULL
);

CREATE TABLE expeditions (
    expeditions_id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    startDate DATE,
    endDate DATE,
    researcher_id INTEGER REFERENCES researcher(researcher_id) NOT NULL
);

CREATE TABLE artifact (
    artifact_id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(100) NOT NULL,
    age INTEGER,
    location VARCHAR(100) NOT NULL,
    researcher_id INTEGER REFERENCES researcher(researcher_id) NOT NULL
);