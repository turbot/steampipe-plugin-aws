select akas, tags_src, title
from aws.aws_ssm_document
where name = '{{ resourceName }}';
