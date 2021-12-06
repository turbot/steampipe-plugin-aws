select arn, author, description, name
from aws_serverlessapplicationrepository_application
where arn = '{{ output.arn.value }}';
