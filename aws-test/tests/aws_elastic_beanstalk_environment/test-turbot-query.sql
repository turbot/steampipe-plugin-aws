select account_id, akas, region, title
from aws.aws_elastic_beanstalk_environment
where environment_name = '{{ resourceName }}';