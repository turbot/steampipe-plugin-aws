select
  db_instance_identifier,
  db_instance_arn,
  db_instance_class,
  dbi_resource_id,
  copy_tags_to_snapshot,
  db_subnet_group_name,
  endpoint_port,
  engine,
  storage_encrypted,
  tags_src
from
  aws.aws_docdb_cluster_instance
WHERE
  db_instance_identifier = '{{ resourceName }}'
