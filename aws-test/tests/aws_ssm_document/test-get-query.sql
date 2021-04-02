select name, tags, title, akas
from aws.aws_ssm_document
where name = '{{ resourceName }}';
