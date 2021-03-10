select repository_name, repository_arn, image_tag_mutability, image_scanning_configuration, partition, region
from aws.aws_ecr_repository
where akas::text = '["{{ output.resource_aka.value }}"]';
