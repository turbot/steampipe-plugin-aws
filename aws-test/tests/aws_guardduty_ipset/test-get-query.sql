select account_id, format, name, partition, region, title
from aws.aws_guardduty_ipset
where detector_id = '{{ output.detector_id.value }}' and ipset_id = '{{ output.ipset_id.value }}';