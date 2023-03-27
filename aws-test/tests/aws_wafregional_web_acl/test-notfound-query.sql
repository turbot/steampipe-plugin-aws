select akas, name, default_action, tags_src, title
from aws.aws_wafregional_web_acl
where name = '{{ resourceName }}-dummy';
