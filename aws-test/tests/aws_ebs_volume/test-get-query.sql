select volume_id, encrypted, tags_src, attachments, multi_attach_enabled
from aws.aws_ebs_volume
where volume_id = '{{ output.resource_id.value }}'
