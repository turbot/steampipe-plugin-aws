select target_group_name, target_group_arn, target_type, port, vpc_id, protocol, matcher_http_code, healthy_threshold_count, unhealthy_threshold_count, health_check_enabled, health_check_interval_seconds, health_check_path, health_check_port, health_check_protocol, health_check_timeout_seconds
from aws.aws_ec2_target_group
where target_group_arn = '{{ output.resource_aka.value }}'
