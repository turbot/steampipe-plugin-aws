select 
  db_instance_identifier,
  title,
  tags,
  akas
from 
  aws.aws_docdb_cluster_instance
where 
  db_instance_identifier = '{{ resourceName }}';
