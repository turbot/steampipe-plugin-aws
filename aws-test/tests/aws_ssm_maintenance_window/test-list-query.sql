select name, window_id, targets, tasks, tags_src, tags, partition, region
from aws.aws_ssm_maintenance_window
where akas::text = '["{{output.resource_aka.value}}"]'