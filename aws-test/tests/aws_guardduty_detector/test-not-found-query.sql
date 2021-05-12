select detector_id, status, service_role, tags
from aws.aws_guardduty_detector
where akas::text = '["{{ output.resource_aka.value }}"::dummy]';
