select db_cluster_identifier, title, tags, akas
from aws.aws_rds_db_cluster
where db_cluster_identifier = '{{ resourceName }}'
