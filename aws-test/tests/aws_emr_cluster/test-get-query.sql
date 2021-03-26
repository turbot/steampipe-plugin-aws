select id, cluster_arn, name, auto_terminate, ebs_root_volume_size, tags
from aws.aws_emr_cluster
where id = '{{ output.resource_id.value }}';
