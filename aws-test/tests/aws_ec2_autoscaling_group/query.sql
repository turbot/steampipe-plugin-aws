select autoscaling_group_arn, name
from aws.aws_ec2_autoscaling_group
where name = '{{resourceName}}'
