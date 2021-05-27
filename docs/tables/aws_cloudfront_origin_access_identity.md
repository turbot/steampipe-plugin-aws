# Table: aws_cloudfront_origin_access_identity

The CloudFront Origin Access Identities page lists of all Origin Access Identities that were created by the RightScale account. An Origin Access Identity (OAI) is used for sharing private content via CloudFront. The OAI is a virtual user identity that will be used to give your CloudFront distribution permission to fetch a private object from your origin server.

## Examples

### Basic Info

```sql
select
  id,
  arn,
  comment,
  s3_canonical_user_id,
  etag
from
  aws_cloudfront_origin_access_identity;
```


### List origin access identity with comment

```sql
select
  id,
  arn,
  comment,
  caller_reference
from
  aws_cloudfront_origin_access_identity
where
  comment <> '';
```