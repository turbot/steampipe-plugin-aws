SELECT
  db_proxy_arn,
  db_proxy_name,
  debug_logging,
  engine_family,
  idle_client_timeout,
  require_tls,
  status,
  tags,
  title
FROM
  aws.aws_rds_db_proxy
WHERE
  db_proxy_name = '{{ resourceName }}'
