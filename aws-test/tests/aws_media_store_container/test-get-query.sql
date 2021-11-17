select arn, name, endpoint, tags
from aws.aws_media_store_container
where name = '{{ resourceName }}';
