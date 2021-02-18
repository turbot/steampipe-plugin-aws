SELECT
  name,
  tags_src
FROM
  aws.aws_rds_db_option_group
WHERE
  name = '{{ resourceName }}'
