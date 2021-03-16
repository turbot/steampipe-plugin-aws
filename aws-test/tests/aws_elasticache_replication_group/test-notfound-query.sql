select title, akas, region, account_id
from aws.aws_elasticache_replication_group
where replication_group_id = 'dummy';