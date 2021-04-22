select account_id, akas, region, title
from aws.aws_ecrpublic_repository
where repository_name = '{{ output.resource_name.value }}';
