select cluster_subnet_group_name, description, akas, tags, title
from aws.aws_redshift_subnet_group
where cluster_subnet_group_name = '{{ resourceName }}';