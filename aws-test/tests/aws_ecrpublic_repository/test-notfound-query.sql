select registry_id, repository_name, arn, region, account_id
from aws.aws_ecrpublic_repository
where repository_name = '{{ resourceName }}.1';

