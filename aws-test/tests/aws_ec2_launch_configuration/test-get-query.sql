select name, launch_configuration_arn, ebs_optimized, image_id, instance_monitoring_enabled, instance_type
from aws.aws_ec2_launch_configuration
where name = '{{resourceName}}'
