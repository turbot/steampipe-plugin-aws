select db_snapshot_identifier, arn, type, allocated_storage, db_instance_identifier, encrypted, engine, engine_version, iam_database_authentication_enabled, license_model, master_user_name, port, storage_type, vpc_id, tag_list
from aws.aws_rds_db_snapshot
where db_snapshot_identifier = '{{ resourceName }}'
