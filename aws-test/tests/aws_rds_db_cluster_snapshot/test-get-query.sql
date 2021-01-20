select db_cluster_snapshot_identifier, arn, type, db_cluster_identifier, engine, engine_version, iam_database_authentication_enabled, license_model, master_user_name, port, storage_encrypted, vpc_id, tag_list
from aws.aws_rds_db_cluster_snapshot
where db_cluster_snapshot_identifier = '{{ resourceName }}'