select endpoint_identifier, arn, endpoint_type
from aws_dms_endpoint
where arn = '{{ output.resource_aka.value }}'
