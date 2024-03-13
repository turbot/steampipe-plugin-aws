SELECT
  db_cluster_snapshot_identifier,
  arn,
  TYPE,
  db_cluster_identifier,
  iam_database_authentication_enabled,
  license_model,
  master_user_name,
  port,
  storage_encrypted,
  vpc_id,
  tags_src
FROM
  aws.aws_rds_db_cluster_snapshot
WHERE
  db_cluster_snapshot_identifier = '{{ resourceName }}'
