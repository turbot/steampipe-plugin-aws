select title, akas, tags, region, account_id
from aws.aws_route53_resolver_rule
where id = '{{ output.resource_id.value }}::asd';
