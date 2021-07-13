select repository_name, description, repository_id, arn
from aws.aws_codecommit_repository
where repository_name = '{{ resourceName }}';