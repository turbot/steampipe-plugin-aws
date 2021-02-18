SELECT
  name,
  parameters,
  tags_src
FROM
  aws.aws_rds_db_parameter_group
WHERE
  name = '{{ resourceName }}'
