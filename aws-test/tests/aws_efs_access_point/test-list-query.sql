select name, access_point_id, access_point_arn
from aws.aws_efs_access_point
where akas::text = '["{{ output.resource_aka.value }}"]';