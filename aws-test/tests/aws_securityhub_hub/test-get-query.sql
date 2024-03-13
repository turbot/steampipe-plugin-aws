select hub_arn, akas, region
from aws_securityhub_hub
where hub_arn = 'arn:aws:securityhub:{{ output.aws_region.value }}:{{ output.aws_account.value }}:hub/default';
