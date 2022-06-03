select name,status, akas
from aws_sfn_state_machine
where arn = '{{ output.resource_aka.value }}';
