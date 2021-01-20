select tags_raw
from aws.aws_sns_topic
where topic_arn = '{{output.resource_aka.value}}'