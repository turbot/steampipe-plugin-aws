select title, akas, region, account_id
from aws.aws_redshift_subnet_group
where cluster_subnet_group_name = 'dummy';