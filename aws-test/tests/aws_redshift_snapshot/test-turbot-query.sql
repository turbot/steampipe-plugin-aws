select title, akas, tags, region, account_id
from aws.aws_redshift_snapshot
where snapshot_identifier = '{{ output.resource_name.value }}';
