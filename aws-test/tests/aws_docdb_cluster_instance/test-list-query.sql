select
  db_instance_identifier,
  db_instance_arn,
  db_instance_class,
  copy_tags_to_snapshot,
  dbi_resource_id,
  db_subnet_group_name,
  endpoint_port,
  engine,
  storage_encrypted,
  tags_src
from
  aws.aws_docdb_cluster_instance
where db_instance_arn = '{{ output.resource_aka.value }}'
