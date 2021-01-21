select snapshot_id, description, volume_id, volume_size, encrypted, owner_id, tags_src
from aws.aws_ebs_snapshot
where snapshot_id = '{{ output.snapshot_id.value }}'
