select name, advanced_event_selectors
from aws.aws_cloudtrail_event_data_store
where name = '{{ resourceName }}';
