select name, tags, title, akas
from aws.aws_ssm_document
where arn = '{{ output.resource_aka.value }}';
