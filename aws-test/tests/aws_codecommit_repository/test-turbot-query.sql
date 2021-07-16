select akas, title
from aws.aws_codecommit_repository
where repository_name = '{{ resourceName }}';