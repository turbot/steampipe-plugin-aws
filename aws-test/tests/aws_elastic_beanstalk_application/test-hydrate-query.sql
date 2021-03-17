select arn , name , tags_src
from aws.aws_elastic_beanstalk_application
where arn = '{{ output.resource_aka.value }}';