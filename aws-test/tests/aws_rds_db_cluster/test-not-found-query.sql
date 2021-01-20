select db_cluster_identifier, arn, status, resource_id
from aws.aws_rds_db_cluster
where db_cluster_identifier = 'dummy-{{ resourceName }}'