select application_id, author, description, name
from aws_serverlessapplicationrepository_application
where name = '{{ output.name.value }}';
