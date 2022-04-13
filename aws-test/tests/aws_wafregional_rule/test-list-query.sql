select name
from aws_wafregional_rule
where akas::text = '["{{ output.resource_aka.value }}"]';