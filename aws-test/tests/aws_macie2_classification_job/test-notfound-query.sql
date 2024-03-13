select region, account_id
from aws.aws_macie2_classification_job
where job_id = '{{output.resource_id.value}}3p';