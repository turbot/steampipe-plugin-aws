select opt_in_status, opt_in_type, opt_in_level, notification_arn, exclusion_list
from aws_servicequotas_auto_management_configuration
where region = '{{ output.resource_aka.value | split:':' | at:3 }}';
