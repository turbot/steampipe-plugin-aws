select registry_id, repository_name, arn, region, account_id
from aws.aws_ecr_repository
where repository_name = '{{ output.resourceName }}.1';

