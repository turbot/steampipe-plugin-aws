select account_id, akas, domain_name, partition, region
from aws_elasticsearch_domain
where akas::text = '["{{ output.resource_aka.value }}"]';