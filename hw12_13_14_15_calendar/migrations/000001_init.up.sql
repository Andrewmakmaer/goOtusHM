CREATE TABLE events (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    starttime TIMESTAMPTZ NOT NULL,
    endtime TIMESTAMPTZ NOT NULL,
    userid VARCHAR(255) NOT NULL,
    callduration TEXT
);