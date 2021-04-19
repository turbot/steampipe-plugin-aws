select account_id, akas, domain_name, ebs_enabled, elasticsearch_version, partition, region, tags, volume_size
from aws.aws_elasticsearch_domain
where domain_name = '{{ resourceName }}';