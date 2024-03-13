select name, partition, region, tags
from aws.aws_macie2_classification_job
where job_id = '{{ output.resource_id.value }}';