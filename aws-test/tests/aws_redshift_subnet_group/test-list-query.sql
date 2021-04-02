select akas, cluster_subnet_group_name, description, subnet_group_status, title, tags
from aws.aws_redshift_subnet_group
where akas::text = '["{{ output.resource_aka.value }}"]';