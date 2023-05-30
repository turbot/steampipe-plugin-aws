select 
  db_instance_identifier,
  db_instance_arn,
  db_instance_class
from 
  aws.aws_docdb_cluster_instance
where 
  db_instance_identifier = 'dummy-{{ resourceName }}'
