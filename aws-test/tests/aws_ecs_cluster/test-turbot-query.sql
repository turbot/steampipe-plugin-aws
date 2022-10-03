select cluster_name, cluster_arn, title, account_id, region, akas
from aws_ecs_cluster
where cluster_name = '{{ resourceName }}';