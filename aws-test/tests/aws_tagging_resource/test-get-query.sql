select name, arn, tags, akas, partition, region, account_id
from aws_tagging_resource
where arn = '{{ output.resource_aka.value }}';
