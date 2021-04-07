select title, akas, tags, region, account_id
from aws.aws_route53_resolver_endpoint
where name = '{{ resourceName }}';
