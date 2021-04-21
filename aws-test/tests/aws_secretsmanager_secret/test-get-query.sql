select name, arn, deleted_date, description, kms_key_id, last_accessed_date, last_rotated_date, owning_service, primary_region, replication_status, rotation_enabled,rotation_lambda_arn, rotation_rules, secret_versions_to_stages, tags_src, tags, title, akas, partition, region, account_id
from aws_secrets_manager_secret
where arn = '{{ output.resource_aka.value }}';