select threat_intel_set_id, format
from aws.aws_guardduty_threat_intel_set
where akas::text = '["{{ output.resource_aka.value }}"::dummy]';
