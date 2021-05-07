select name, arn, assessment_target_arn, duration_in_seconds
from aws.aws_inspector_assessment_template
where akas::text = '["{{ output.resource_aka.value }}"]';