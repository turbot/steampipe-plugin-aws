select name, arn
from aws.aws_ec2_application_load_balancer
where name = '{{resourceName}}'
