select akas, name, scope, description, visibility_config, default_action, tags_src, title
from aws.aws_wafv2_web_acl
where name = '{{ output.resource_name.value }}';