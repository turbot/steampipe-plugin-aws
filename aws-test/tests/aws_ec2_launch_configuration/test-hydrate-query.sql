select name, akas, title
from aws.aws_ec2_launch_configuration
where name = '{{resourceName}}'
