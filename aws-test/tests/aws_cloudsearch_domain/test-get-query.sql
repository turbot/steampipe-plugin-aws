select domain_name, domain_id, arn
from aws.aws_cloudsearch_domain
where domain_name = '{{ resourceName }}';
