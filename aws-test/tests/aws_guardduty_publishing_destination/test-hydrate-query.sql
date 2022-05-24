select destination_id, arn, detector_id, partition, region, title
from aws_guardduty_publishing_destination
where akas::text = '["{{ output.resource_aka.value }}"]';
