select repository_name, akas, tags, title
from aws.aws_ecr_repository
where repository_name = '{{ output.resource_name.value }}'
