select akas, tags, title
from aws.aws_ec2_gateway_load_balancer
where name = '{{ resourceName }}';
