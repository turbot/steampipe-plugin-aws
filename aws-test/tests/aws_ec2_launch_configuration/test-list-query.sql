select name, launch_configuration_arn
from aws.aws_ec2_launch_configuration
where name = '{{resourceName}}'
