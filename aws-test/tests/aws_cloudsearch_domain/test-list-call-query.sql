select arn, domain_name
from aws.aws_cloudsearch_domain
where akas::text = '["{{ output.resource_aka.value }}"]';
