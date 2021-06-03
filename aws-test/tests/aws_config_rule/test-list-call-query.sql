select name, rule_id, arn, title 
from aws.aws_config_rule
where akas::text = '["{{ output.resource_aka.value }}"]'