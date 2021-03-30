select title, akas, tags, region, account_id
from aws.aws_ssm_document
where name = '{{ resourceName }}::xzq';
