select arn, name
from aws.aws_cloudtrail_trail
where akas::text = '["{{ output.resource_aka.value }}"]';
