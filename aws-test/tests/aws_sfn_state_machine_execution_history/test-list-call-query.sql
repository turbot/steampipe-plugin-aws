select id, execution_arn, type, akas
from aws_sfn_state_machine_execution_history
where id = '{{ output.id.value }}' and execution_arn = '{{ output.execution_arn.value }}';
