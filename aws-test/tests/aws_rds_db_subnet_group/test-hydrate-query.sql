SELECT
  name,
  tags_src
FROM
  aws.aws_rds_db_subnet_group
WHERE
  name = '{{ resourceName }}'
