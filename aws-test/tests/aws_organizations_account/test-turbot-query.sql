select title, akas, partition, region, account_id
from aws.aws_organizations_account
where id = '{{ output.resource_id.value }}';
