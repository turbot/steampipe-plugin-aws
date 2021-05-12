select workload_id, workload_name, title, akas, account_id
from aws.aws_wellarchitected_workload
where workload_id = '{{ output.resource_id.value }}';
