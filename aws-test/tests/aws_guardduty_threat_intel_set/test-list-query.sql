select account_id, format, name, partition, region, title
from aws.aws_guardduty_threat_intel_set
where akas::text = '["{{ output.resource_aka.value }}"]';