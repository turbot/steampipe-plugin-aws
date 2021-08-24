select title, akas, tags, region, account_id
from aws_tagging_resource
where arn = '{{ output.resource_aka.value }}';