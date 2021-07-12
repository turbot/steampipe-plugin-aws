select repository_name, repository_id, arn, description
from aws.aws_codecommit_repository
where repository_name = '{{ resourceName }}';