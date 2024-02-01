select name
from aws.aws_cloudwatch_log_group
where name = '{{ resourceName }}' and region = '{{ output.region_name.value }}'
