select block_public_acls, ignore_public_acls, bucket_policy_is_public, logging, name, partition, restrict_public_buckets, tags_src, akas, tags, title, versioning_enabled, versioning_mfa_delete
from aws.aws_s3_bucket
where akas::text = '["arn:aws:s3:::{{ resourceName }}"]';
