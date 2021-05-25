select tags_src, application_description, service_execution_role, cloud_watch_logging_option_descriptions
from aws.aws_kinesisanalyticsv2_application
where application_name = '{{ resourceName }}';
