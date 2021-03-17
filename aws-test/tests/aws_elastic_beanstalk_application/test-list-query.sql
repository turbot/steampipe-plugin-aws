select arn , name , partition , region  , tags , title
from aws.aws_elastic_beanstalk_application
where akas::text = '["{{ output.resource_aka.value }}"]';