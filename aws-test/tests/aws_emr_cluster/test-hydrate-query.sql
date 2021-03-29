select id, cluster_arn, name, auto_terminate, tags, akas, region, account_id
from aws.aws_emr_cluster
where id = '{{ output.resource_id.value }}';