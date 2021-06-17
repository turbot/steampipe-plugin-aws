select application_name, application_arn, application_status, application_version_id
from aws.aws_kinesisanalyticsv2_application
where akas::text = '["{{ output.resource_aka.value }}"]';