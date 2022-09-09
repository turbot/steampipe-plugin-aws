select db_cluster_identifier, arn, status
from aws.aws_docdb_cluster
where db_cluster_identifier = 'dummy-{{ resourceName }}'
