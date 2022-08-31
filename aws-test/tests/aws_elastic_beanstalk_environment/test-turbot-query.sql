select account_id, akas, region, title
from aws_elastic_beanstalk_environment
where environment_name = '{{ resourceName }}';