select akas, environment_name, tags, title
from aws.aws_elastic_beanstalk_environment
where environment_name = '{{ resourceName }}';