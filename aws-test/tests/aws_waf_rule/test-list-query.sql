select name, partition, tags, title
from aws.aws_waf_rule
where akas::text = '["{{ output.resource_aka.value }}"]';