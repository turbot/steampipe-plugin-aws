select target_group_arn, target_health_descriptions, tags_src, akas, tags, title
from aws.aws_ec2_target_group
where target_group_arn = '{{ output.resource_aka.value }}'
