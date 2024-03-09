CREATE DATABASE goProject;

CREATE TABLE researchers (
    researchers_id INTEGER PRIMARY KEY,
    name VARCHAR(25) NOT NULL,
    specialization VARCHAR(25) NOT NULL,
    project VARCHAR(25) NOT NULL
);

CREATE TABLE expeditions (
    expeditions_id INTEGER PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    startDate DATE,
    endDate DATE,
    researchers_id INTEGER REFERENCES researchers(researchers_id) NOT NULL
);


CREATE TABLE artifact (
    artifact_id INTEGER PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    description VARCHAR(100) NOT NULL,
    age INTEGER,
    location VARCHAR(50) NOT NULL,
    researchers_id INTEGER REFERENCES researchers(researchers_id) NOT NULL
);

INSERT INTO researchers (researchers_id, name, specialization, project) VALUES
    (1, 'Indiana Jones', 'Archaeology', 'none'),
    (2, 'Jane Goodall', 'Anthropology', 'none'),
     (3, 'Alan Grant', 'Paleontology', 'none');

INSERT INTO expeditions (expeditions_id, title, startDate, endDate, researchers_id) VALUES
    (10, 'Archaeological Expedition 1', '2022-01-01', '2022-02-01', 1),
    (20, 'Anthropological Expedition 1', '1980-03-01', '1990-04-01', 2),
    (30, 'Paleontological Expedition 1', '1990-05-01', '1992-06-01', 3);

INSERT INTO artifact (artifact_id, title, description, age, location, researchers_id) VALUES
    (101, 'Ancient Relic', 'A mysterious artifact from the past', 500, 'Lost City', 1),
    (102, 'Rare Fossil', 'Well-preserved dinosaur fossil', 65, 'Dinosaur Valley', 2),
    (103, 'Enigmatic Pottery', 'Unusual pottery with unknown origin', 100, 'Mystic Ruins', 3);