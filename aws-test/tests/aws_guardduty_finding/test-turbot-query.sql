select account_id, akas, title
from aws.aws_guardduty_finding
where akas::text = '["{{ output.resource_aka.value }}"]';