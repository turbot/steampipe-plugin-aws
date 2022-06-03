select akas
from aws_sfn_state_machine_execution
where state_machine_arn = '{{ output.state_machine_arn.value }}';
