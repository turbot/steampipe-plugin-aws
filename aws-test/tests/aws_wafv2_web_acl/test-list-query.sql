select akas, name, scope, description, title
from aws.aws_wafv2_web_acl
where akas::text = '["{{ output.resource_aka.value }}"]';
