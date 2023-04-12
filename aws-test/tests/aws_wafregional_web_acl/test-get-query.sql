select akas, name, default_action, tags_src, title
from aws.aws_wafregional_web_acl
where web_acl_id = '{{ output.resource_id.value }}';