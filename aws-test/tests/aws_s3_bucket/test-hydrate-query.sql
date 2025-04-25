select name, logging, acl, replication, versioning_enabled, versioning_mfa_delete, bucket_policy_is_public, object_lock_configuration, region
from aws.aws_s3_bucket
where name = '{{ resourceName }}';
