CREATE TABLE IF NOT EXISTS goods
(
    Id         INT,
    ProjectId INT,
    Name       String,
    Description Nullable(String),
    Priority   INT,
    Removed    BOOLEAN  DEFAULT false,
    EventTime DateTime DEFAULT now()

) ENGINE = MergeTree()
      ORDER BY (Id, ProjectId, Name);
