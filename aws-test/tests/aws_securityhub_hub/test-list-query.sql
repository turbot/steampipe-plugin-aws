select hub_arn, auto_enable_controls, akas, region
from aws_securityhub_hub
where akas = '["arn:aws:securityhub:{{ output.aws_region.value }}:{{ output.aws_account.value }}:hub/default"]';
