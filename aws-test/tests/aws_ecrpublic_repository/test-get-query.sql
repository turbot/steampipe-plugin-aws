select akas, partition, region, repository_name
from aws.aws_ecrpublic_repository
where repository_name = '{{ output.resource_name.value }}';
