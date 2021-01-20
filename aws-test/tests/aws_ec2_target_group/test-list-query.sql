select target_group_name, target_group_arn
from aws.aws_ec2_target_group
where target_group_name = '{{resourceName}}'
