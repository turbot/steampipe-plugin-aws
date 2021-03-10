select akas, image_tag_mutability,partition, region, repository_arn, repository_name 
from aws.aws_ecr_repository
where repository_name = '{{ output.resource_name.value }}'
