select account_id, akas, title
from aws.aws_guardduty_threat_intel_set
where akas::text = '["{{ output.resource_aka.value }}"]';
