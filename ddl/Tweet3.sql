CREATE TABLE Tweet3 (
  Id STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL,
  CommitedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (Id);