select akas , name , region , tags , title
from aws.aws_redshift_parameter_group
where name = '{{ resourceName }}'; 