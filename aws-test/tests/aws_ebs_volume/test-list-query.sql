select tags_src, volume_id, title, tags, akas
from aws.aws_ebs_volume
where akas::text = '["{{ output.resource_aka.value }}"]'
