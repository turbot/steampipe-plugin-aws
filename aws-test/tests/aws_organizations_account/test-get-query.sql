select id, arn, email, name, status
from aws.aws_organizations_account
where id = '{{ output.resource_id.value }}';
