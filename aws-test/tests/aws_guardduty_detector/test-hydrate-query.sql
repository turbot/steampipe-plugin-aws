select detector_id, status, tags, akas, title
from aws.aws_guardduty_detector
where akas::text = '["{{ output.resource_aka.value }}"]';
