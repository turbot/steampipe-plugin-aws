select akas, tags, title
from aws.aws_inspector_assessment_template
where arn = '{{ output.resource_aka.value }}';