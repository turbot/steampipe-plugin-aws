select akas, title
from aws.aws_opensearch_domain
where domain_name = '{{ resourceName }}';