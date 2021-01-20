select name, title, tags, akas
from aws.aws_ec2_classic_load_balancer
where name = '{{ resourceName }}'