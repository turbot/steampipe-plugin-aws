select akas, title
from aws.aws_inspector_assessment_target
where arn = '{{ output.resource_aka.value }}';