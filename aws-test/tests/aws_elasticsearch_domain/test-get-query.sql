select account_id, akas, domain_name, ebs_options, elasticsearch_version, partition, region, snapshot_options
from aws.aws_elasticsearch_domain
where domain_name = '{{ resourceName }}';