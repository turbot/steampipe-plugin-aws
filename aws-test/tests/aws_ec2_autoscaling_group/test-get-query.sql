select autoscaling_group_arn, name, availability_zones, launch_configuration_name, tags_src
from aws.aws_ec2_autoscaling_group
where name = '{{ resourceName }}'
