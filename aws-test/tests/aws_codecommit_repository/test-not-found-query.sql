select name
from aws.aws_codecommit_repository
where name = 'dummy{{ resourceName }}';