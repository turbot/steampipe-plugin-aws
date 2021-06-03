select account_id, region, title, aka
from aws.aws_macie2_classification_job
where job_id = '{{output.resource_id.value}}';
