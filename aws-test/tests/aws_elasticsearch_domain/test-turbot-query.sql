select akas, title
from aws.aws_elasticsearch_domain
where domain_name = '{{ resourceName }}';