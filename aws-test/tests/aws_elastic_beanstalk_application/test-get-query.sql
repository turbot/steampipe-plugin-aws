select arn, name, tags
from aws.aws_elastic_beanstalk_application
where name = '{{ resourceName }}';