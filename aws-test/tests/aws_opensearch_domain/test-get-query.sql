select account_id, akas, domain_name, engine_version, partition, region
from aws.aws_opensearch_domain
where domain_name = '{{ resourceName }}';