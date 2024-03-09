project
# Golang_2024project
# Artifacts

## Members
```
Orazbai Dariga 22B030617
Kemelbay Merey 22B030375
```

## Artifacts REST API
```
POST /artifacts
GET /artifacts/id
PUT /artifacts/id
DELETE /artifacts/
```
## DB Structure

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