select account_id, akas, domain_name, ebs_options, elasticsearch_version, partition, region, snapshot_options, tags
from aws.aws_elasticsearch_domain
where arn = '{{ output.resource_aka.value }}';