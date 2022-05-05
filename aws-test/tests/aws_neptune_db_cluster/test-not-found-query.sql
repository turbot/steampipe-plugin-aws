select db_cluster_identifier, arn, status
from aws_neptune_db_cluster
where db_cluster_identifier = 'dummy-{{ resourceName }}';
