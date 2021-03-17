select title, akas, region, account_id
from aws.aws_elasticache_subnet_group
where cache_subnet_group_name = 'dummy';