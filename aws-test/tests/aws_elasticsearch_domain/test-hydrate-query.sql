select account_id, akas, domain_name, elasticsearch_version, partition, region, snapshot_options
from aws_elasticsearch_domain
where arn = '{{ output.resource_aka.value }}';