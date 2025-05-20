select application_name, application_version_id, runtime_environment, service_execution_role, title, akas, partition, account_id
from aws.aws_kinesisanalyticsv2_application
where application_name = '{{ resourceName }}';
