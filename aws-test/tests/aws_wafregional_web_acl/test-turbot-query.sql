select akas, tags, title
from aws.aws_wafregional_web_acl
where name = '{{ resourceName }}';