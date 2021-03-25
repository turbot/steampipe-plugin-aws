select account_id, akas, region, title
from aws.aws_ssm_association
where association_id = '{{ output.resource_id.value }}';
