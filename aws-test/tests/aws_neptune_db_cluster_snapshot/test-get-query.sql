SELECT
  db_cluster_snapshot_identifier,
  db_cluster_snapshot_arn,
  db_cluster_identifier,
  snapshot_type,
  iam_database_authentication_enabled,
  license_model,
  port,
  storage_encrypted
FROM
  aws.aws_neptune_db_cluster_snapshot
WHERE
  db_cluster_snapshot_identifier = '{{ resourceName }}'
