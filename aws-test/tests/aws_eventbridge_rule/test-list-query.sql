select arn , event_bus_name , name , partition , region , state , tags , title
from aws.aws_eventbridge_rule
where akas::text = '["{{ output.resource_aka.value }}"]';