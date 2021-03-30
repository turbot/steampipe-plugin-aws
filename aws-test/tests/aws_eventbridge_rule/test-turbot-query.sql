select akas , name , region , tags , title
from aws.aws_eventbridge_rule
where name = '{{ resourceName }}';