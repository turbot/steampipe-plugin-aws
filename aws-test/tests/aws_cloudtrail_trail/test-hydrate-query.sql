select tags_src, event_selectors, is_logging
from aws.aws_cloudtrail_trail
where name = '{{ resourceName }}';
