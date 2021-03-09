select name, arn, log_file_validation_enabled, has_custom_event_selectors, has_insight_selectors, home_region, include_global_service_events, is_multi_region_trail, is_organization_trail, s3_bucket_name, s3_key_prefix
from aws.aws_cloudtrail_trail
where name = '{{ resourceName }}';
