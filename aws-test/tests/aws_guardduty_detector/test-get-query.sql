select detector_id, status, finding_publishing_frequency, tags
from aws.aws_guardduty_detector
where detector_id = '{{ output.resource_id.value }}';