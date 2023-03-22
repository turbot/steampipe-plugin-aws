select db_cluster_identifier, arn, db_cluster_resource_id, status, copy_tags_to_snapshot, cross_account_clone, db_subnet_group, deletion_protection, endpoint, engine, iam_database_authentication_enabled, multi_az, port, storage_encrypted, tags_src
from aws_neptune_db_cluster
where db_cluster_identifier = '{{ resourceName }}';
