select metric_name, partition, title
from aws.aws_waf_rate_based_rule
where akas::text = '["{{ output.resource_aka.value }}"]';