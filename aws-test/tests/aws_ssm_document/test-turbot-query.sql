select akas, name, region, tags, title
from aws.aws_ssm_document
where name = '{{ resourceName }}';
