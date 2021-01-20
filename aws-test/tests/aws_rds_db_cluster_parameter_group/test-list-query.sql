select name, arn
from aws.aws_rds_db_cluster_parameter_group
where arn = '{{ output.resource_aka.value }}'