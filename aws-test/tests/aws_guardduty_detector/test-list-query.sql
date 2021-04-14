select detector_id, status, finding_publishing_frequency
from aws.aws_guardduty_detector
where akas::text = '["{{ output.resource_aka.value }}"]';
