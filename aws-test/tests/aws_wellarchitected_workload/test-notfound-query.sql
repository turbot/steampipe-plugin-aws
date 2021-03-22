select workload_id,  workload_name
from aws.aws_wellarchitected_workload
where workload_name = '{{ output.resource_name.value }}::dsf';