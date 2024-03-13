select akas, title
from aws.aws_cloudsearch_domain
where domain_name = '{{ resourceName }}';
