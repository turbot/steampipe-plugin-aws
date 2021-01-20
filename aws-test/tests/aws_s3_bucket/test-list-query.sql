select block_public_acls, block_public_policy, ignore_public_acls, bucket_policy_is_public, lifecycle_rules, logging, name, partition, restrict_public_buckets, server_side_encryption_configuration, tags_raw, akas, tags, title, versioning_enabled, versioning_mfa_delete
from aws.aws_s3_bucket
where akas::text = '["arn:aws:s3:::{{ resourceName }}"]'
