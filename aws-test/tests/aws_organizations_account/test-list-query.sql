select id, arn, email, name, status
from aws.aws_organizations_account
where arn = '{{ output.resource_arn.value }}';
