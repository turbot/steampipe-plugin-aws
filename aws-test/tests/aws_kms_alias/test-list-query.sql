
select arn
from aws.aws_kms_alias
where akas::text = '["{{ output.resource_aka.value }}"]'
