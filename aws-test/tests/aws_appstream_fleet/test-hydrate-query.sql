select name, akas, tags, title
from aws_appstream_fleet
where arn = '{{ output.resource_aka.value }}'
