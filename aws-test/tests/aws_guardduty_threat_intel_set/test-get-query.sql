select account_id, format, name, partition, region, title
from aws.aws_guardduty_threat_intel_set
where detector_id = '{{ output.detector_id.value }}' and threat_intel_set_id = '{{ output.resource_id.value }}';
