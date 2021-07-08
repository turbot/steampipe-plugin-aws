select job_id, name, arn, partition, region
from aws.aws_macie2_classification_job
where name = '{{ output.resource_name.value }}';