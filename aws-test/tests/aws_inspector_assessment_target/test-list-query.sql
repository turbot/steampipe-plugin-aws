select name, arn, resource_group_arn
from aws.aws_inspector_assessment_target
where akas::text = '["{{ output.resource_aka.value }}"]';