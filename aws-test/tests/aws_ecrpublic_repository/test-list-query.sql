select repository_name, arn, partition, region
from aws.aws_ecrpublic_repository
where akas::text = '["{{ output.resource_aka.value }}"]';
