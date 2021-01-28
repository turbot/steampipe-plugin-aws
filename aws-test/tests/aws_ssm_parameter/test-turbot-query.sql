select title, akas, tags, region, account_id
from aws.aws_ssm_parameter
where name = '{{ resourceName }}'
