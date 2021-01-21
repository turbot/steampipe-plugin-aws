select access_log_emit_interval, access_log_enabled, access_log_s3_bucket_name, access_log_s3_bucket_prefix, additional_attributes, connection_draining_enabled, connection_draining_timeout, connection_settings_idle_timeout, cross_zone_load_balancing_enabled, tags_src
from aws.aws_ec2_classic_load_balancer
where name = '{{ resourceName }}'
