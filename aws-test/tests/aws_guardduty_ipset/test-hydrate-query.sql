select account_id, format, name, partition, region, title
from aws.aws_guardduty_ipset
where akas::text = '["{{ output.resource_aka.value }}"]';