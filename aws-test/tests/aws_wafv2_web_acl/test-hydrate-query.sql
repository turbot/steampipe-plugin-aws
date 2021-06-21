select name, scope, akas, tags, title, logging_configuration
from aws.aws_wafv2_web_acl
where id = '{{ output.resource_id.value }}';
