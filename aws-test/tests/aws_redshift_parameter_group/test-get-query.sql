select family, name , tags
from aws.aws_redshift_parameter_group
where name = '{{ resourceName }}';