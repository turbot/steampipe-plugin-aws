select destination_id, arn
from aws_guardduty_publishing_destination
where akas::text = '["{{ output.resource_aka.value }}"::dummy]';
