
select alias_name
from aws.aws_kms_alias
where alias_name='alias/my-key-alias_{{ resourceName }}'
