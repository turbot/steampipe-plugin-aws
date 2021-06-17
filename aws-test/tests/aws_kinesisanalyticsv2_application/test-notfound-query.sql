select application_name, application_status, application_version_id
from aws.aws_kinesisanalyticsv2_application
where application_name = '{{ resourceName }}NF';