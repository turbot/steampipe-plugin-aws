select tags_src
from aws.aws_sns_topic
where topic_arn = '{{output.resource_aka.value}}'
