select account_id, akas, region, title
from aws.aws_ecr_repository
where repository_name = '{{ output.resource_name.value }}';
