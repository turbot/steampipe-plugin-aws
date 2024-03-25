select akas
from aws.aws_backup_job
where backup_job_id = '{{ output.id.value }}';