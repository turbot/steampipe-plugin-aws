select name, policy, policy_std, logging, acl, lifecycle_rules, server_side_encryption_configuration, replication, versioning_enabled, versioning_mfa_delete, bucket_policy_is_public, object_lock_configuration
from aws.aws_s3_bucket
where name = '{{ resourceName }}'
