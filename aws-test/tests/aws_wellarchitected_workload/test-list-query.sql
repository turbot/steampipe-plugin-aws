select workload_id, workload_name, title, akas
from aws.aws_wellarchitected_workload
where workload_name = '{{ output.resource_name.value }}';
