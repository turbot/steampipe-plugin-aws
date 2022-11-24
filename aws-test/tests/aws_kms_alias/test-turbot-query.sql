select akas, title
from aws.aws_kms_alias
where arn ='{{ output.resource_aka.value }}'
