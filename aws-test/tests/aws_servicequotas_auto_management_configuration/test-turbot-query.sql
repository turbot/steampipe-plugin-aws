select title, akas
from aws_servicequotas_auto_management_configuration
where region = '{{ output.resource_aka.value | split:':' | at:3 }}';
