select name
from aws.aws_emr_security_configuration
where name = 'dummy-{{ resourceName }}';