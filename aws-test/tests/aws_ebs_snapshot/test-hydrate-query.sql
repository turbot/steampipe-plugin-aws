select snapshot_id, create_volume_permissions, akas, tags, title
from aws.aws_ebs_snapshot
where snapshot_id = '{{ output.snapshot_id.value }}'
