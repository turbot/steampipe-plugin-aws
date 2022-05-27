select akas, title
from aws_serverlessapplicationrepository_application
where name = '{{ output.name.value }}';
