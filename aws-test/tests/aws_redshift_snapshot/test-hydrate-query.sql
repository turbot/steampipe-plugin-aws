select snapshot_identifier, akas, tags
from aws.aws_redshift_snapshot
where snapshot_identifier = '{{ output.resource_name.value }}';
