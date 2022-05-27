select akas, name, region, tags, title
from aws.aws_media_store_container
where name = '{{ resourceName }}';
