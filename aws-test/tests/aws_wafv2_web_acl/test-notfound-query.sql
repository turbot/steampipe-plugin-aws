select title, akas, region, account_id
from aws.aws_wafv2_web_acl
where name = '{{ resourceName }}-dummy';