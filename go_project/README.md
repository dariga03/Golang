project
# Golang_2024project
# Archaeological Expeditions

## Members
```
Orazbai Dariga 22B030617
Kemelbay Merey 22B030375
```

## Description of the project
```
Project Archaeological Expeditions {
  database_type: 'PostgreSQL'
  Note: 'The Archaeological Expeditions project is a database organized around three key tables: expeditions, researchers, and artifacts. These tables are interconnected, creating a structure for tracking information about artifacts discovered during archaeological expeditions and related to researchers who participated in these expeditions.'
}
```

## Researcher REST API
```
POST /researchers
GET /researchers/id
PUT /researchers/id
DELETE /researchers/
```

## Expedition REST API
```
POST /expeditions
GET /expeditions/id
PUT /expeditions/id
DELETE /expeditions/
GET /researchers/id/expeditions
```

## Artifact REST API
```
POST /artifacts
GET /artifacts/id
PUT /artifacts/id
DELETE /artifacts/
GET /researchers/id/artifacts
```
## DB Structure
```
TABLE researchers (
    researchers_id INTEGER PRIMARY KEY,
    name VARCHAR(25) NOT NULL,
    specialization VARCHAR(25) NOT NULL,
    project VARCHAR(25) NOT NULL
);



TABLE expeditions (
    expeditions_id INTEGER PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    startDate DATE,
    endDate DATE,
    researchers_id INTEGER REFERENCES researchers(researchers_id) NOT NULL
);


TABLE artifact (
    artifact_id INTEGER PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    description VARCHAR(100) NOT NULL,
    age INTEGER,
    location VARCHAR(50) NOT NULL,
    researchers_id INTEGER REFERENCES researchers(researchers_id) NOT NULL
);
```