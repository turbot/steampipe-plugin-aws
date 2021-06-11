select default_ebs_encryption_enabled, default_ebs_encryption_key, title, region
from aws.aws_ec2_regional_settings
where region = '{{ output.aws_region.value }}'