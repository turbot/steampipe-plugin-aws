select db_cluster_identifier, arn, status, resource_id
from aws_pinpoint_app
where id = 'dummy-{{ resourceName }}';
