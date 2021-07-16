select repository_name, tags
from aws.aws_codecommit_repository
where repository_name = '{{ resourceName }}';