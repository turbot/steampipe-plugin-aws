select repository_name
from aws.aws_codecommit_repository
where repository_name = 'dummy-{{ resourceName }}';