select akas, tags, title
from aws.aws_kinesisanalyticsv2_application
where application_name = '{{ resourceName }}';
