select name, id, comment
from aws.aws_route53_zone
where name = '{{ resourceName }}.com.'