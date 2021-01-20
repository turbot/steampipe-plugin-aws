select snapshot_id, volume_id
from aws.aws_ebs_snapshot
where snapshot_id = '{{ output.snapshot_id.value }}'
