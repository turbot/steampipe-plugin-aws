select account_id, akas, title
from aws_guardduty_publishing_destination
where akas::text = '["{{ output.resource_aka.value }}"]';
