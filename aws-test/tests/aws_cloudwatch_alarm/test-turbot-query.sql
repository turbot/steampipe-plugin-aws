select name, alarm_arn, tags, akas, title, partition, region, account_id
from aws.aws_cloudwatch_alarm
where akas::text = '["{{ output.resource_aka.value }}"]';
