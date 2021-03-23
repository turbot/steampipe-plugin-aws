select akas, title
from aws.aws_redshift_snapshot
where snapshot_identifier = '{{ output.resource_name.value }}';
