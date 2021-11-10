select application_id, author, description
from aws_serverlessapplicationrepository_application
where name = 'dummy-{{ output.name.value }}';
