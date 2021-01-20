select name, title, tags, akas
from aws.aws_rds_db_cluster_parameter_group
where name = '{{ resourceName }}'