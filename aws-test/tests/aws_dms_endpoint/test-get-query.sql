select endpoint_identifier, arn
from aws_dms_endpoint
where arn = '{{ output.resource_aka.value }}';
