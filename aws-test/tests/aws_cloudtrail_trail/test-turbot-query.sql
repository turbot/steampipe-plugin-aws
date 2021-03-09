select tags, akas, title
from aws.aws_cloudtrail_trail
where name = '{{ resourceName }}';
