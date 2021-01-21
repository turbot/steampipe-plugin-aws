
select id
from aws.aws_kms_key
where akas::text = '["{{ output.resource_aka.value }}"]'
