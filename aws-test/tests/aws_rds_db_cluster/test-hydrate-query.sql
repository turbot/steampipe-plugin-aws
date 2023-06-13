select db_cluster_identifier, arn, resource_id, allocated_storage, copy_tags_to_snapshot, cross_account_clone, db_subnet_group, deletion_protection, endpoint, engine, engine_mode, global_write_forwarding_requested, http_endpoint_enabled, iam_database_authentication_enabled, master_user_name, multi_az, port, storage_encrypted, tags_src
from aws.aws_rds_db_cluster
where db_cluster_identifier = '{{ resourceName }}'
