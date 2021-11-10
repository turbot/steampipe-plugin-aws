select application_id, author, description, name
from aws_serverlessapplicationrepository_application
where application_id = '{{ output.application_id.value }}';
