select name, description, id
from aws.aws_codecommit_repository
where name = '{{ resourceName }}';