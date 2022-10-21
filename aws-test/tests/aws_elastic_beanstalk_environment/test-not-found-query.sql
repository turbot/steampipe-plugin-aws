select application_name, environment_id,environment_name, partition, region
from aws_elastic_beanstalk_environment
where environment_name = '{{ resourceName }}.ab';