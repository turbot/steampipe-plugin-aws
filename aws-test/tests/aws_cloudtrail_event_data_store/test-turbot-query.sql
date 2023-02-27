select akas, title
from aws.aws_cloudtrail_event_data_store
where name = '{{ resourceName }}';
