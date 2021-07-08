SELECT
  name,
  tags_src
FROM
  aws.aws_rds_db_cluster_parameter_group
WHERE
  name = '{{ resourceName }}'
