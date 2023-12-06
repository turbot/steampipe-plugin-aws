select thing_type_name, arn
from aws.aws_iot_thing_type
where akas::text = '["{{ output.resource_aka.value }}"]';