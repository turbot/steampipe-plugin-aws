select event_data_store_arn, name
from aws_cloudtrail_event_data_store
where akas::text = '["{{ output.resource_aka.value }}"]';
