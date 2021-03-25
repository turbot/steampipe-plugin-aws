select arn , event_pattern , name , tags
from aws.aws_eventbridge_rule
where name = '{{ resourceName }}';
