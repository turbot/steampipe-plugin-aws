select job_id, name, partition, region
from aws.aws_macie2_classification_job
where name = '{{ output.resource_name.value }}';