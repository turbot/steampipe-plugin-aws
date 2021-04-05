select workload_id, workload_name, title, akas, account_id, region, partition
from aws.aws_wellarchitected_workload
where workload_name = '{{ output.resource_name.value }}';