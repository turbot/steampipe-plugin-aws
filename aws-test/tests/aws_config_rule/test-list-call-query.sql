select name, rule_id, rule_arn, title 
from aws.aws_config_rule
where akas::text = '["{{ output.resource_aka.value }}"]'