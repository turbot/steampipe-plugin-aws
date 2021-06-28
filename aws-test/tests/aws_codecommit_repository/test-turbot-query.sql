select akas, title
from aws.aws_codecommit_repository
where name = '{{ resourceName }}'
