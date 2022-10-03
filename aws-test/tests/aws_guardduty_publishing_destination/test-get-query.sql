select destination_id, arn, detector_id, partition, region
from aws_guardduty_publishing_destination
where detector_id = '{{ output.detector_id.value }}' and destination_id = '{{ output.resource_id.value }}';
