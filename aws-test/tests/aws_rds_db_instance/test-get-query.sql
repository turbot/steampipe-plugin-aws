select db_instance_identifier, arn, resource_id, class, allocated_storage, auto_minor_version_upgrade, copy_tags_to_snapshot, customer_owned_ip_enabled, port, db_subnet_group_name, deletion_protection, endpoint_port, engine, engine_version, iam_database_authentication_enabled, master_user_name, multi_az, performance_insights_enabled, publicly_accessible, storage_encrypted, storage_type, tag_list
from aws.aws_rds_db_instance
where db_instance_identifier = '{{ resourceName }}'
