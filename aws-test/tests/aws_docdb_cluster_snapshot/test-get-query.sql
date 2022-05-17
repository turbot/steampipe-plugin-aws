SELECT
  db_cluster_snapshot_identifier,
  arn,
  type,
  db_cluster_identifier,
  engine,
  storage_encrypted
FROM
  aws.aws_docdb_cluster_snapshot
WHERE
  db_cluster_snapshot_identifier = '{{ resourceName }}'
