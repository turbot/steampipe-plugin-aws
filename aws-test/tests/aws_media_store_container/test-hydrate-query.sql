select arn, name, endpoint, tags_src
from aws.aws_media_store_container
where name = '{{ resourceName }}';
