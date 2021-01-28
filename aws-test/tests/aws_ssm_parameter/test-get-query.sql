select name, tags, title, akas
from aws.aws_ssm_parameter
where name = '{{ resourceName }}'
