select autoscaling_group_arn, name, availability_zones, launch_configuration_name, tags_raw, tags, akas, title
from aws.aws_ec2_autoscaling_group
where akas::text = '["{{ output.resource_aka.value }}"]'