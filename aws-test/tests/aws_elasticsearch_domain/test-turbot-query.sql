select akas, title
from aws_elasticsearch_domain
where domain_name = '{{ resourceName }}';