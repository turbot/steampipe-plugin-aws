select arn , name , tags_src
from aws.aws_eventbridge_rule
where arn = '{{ output.resource_aka.value }}';