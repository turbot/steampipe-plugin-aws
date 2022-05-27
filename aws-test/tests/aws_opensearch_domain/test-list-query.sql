select account_id, akas, domain_name, partition, region
from aws.aws_opensearch_domain
where akas::text = '["{{ output.resource_aka.value }}"]';