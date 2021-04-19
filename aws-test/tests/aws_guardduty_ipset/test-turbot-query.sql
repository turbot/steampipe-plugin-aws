select account_id, akas, title
from aws.aws_guardduty_ipset
where akas::text = '["{{ output.resource_aka.value }}"]';