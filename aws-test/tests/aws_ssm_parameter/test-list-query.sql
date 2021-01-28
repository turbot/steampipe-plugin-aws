select name, data_type, key_id, policies, tags_src, tags, partition, region, tier
from aws.aws_ssm_parameter
where akas::text = '["{{output.resource_aka.value}}"]'
