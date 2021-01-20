select tags, title, akas, region
from aws.aws_iam_user
where name = '{{ resourceName }}'
