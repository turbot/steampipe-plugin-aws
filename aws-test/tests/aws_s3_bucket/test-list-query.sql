select bucket_policy_is_public, logging, name, partition, tags_src, akas, tags, title, versioning_enabled, versioning_mfa_delete
from aws.aws_s3_bucket
where akas::text = '["arn:aws:s3:::{{ resourceName }}"]';
