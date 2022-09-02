select akas, tags, title
from aws.aws_waf_web_acl
where name = '{{ resourceName }}';