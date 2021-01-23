select name, policies
from aws.aws_ec2_autoscaling_group
where name = '{{resourceName}}'