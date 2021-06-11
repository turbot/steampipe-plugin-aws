# Table: aws_cloudfront_origin_access_identity

An origin access identity is a special CloudFront user that is associated with a distribution. You associate the origin access identity with origins, so that you can secure all or just some of your Amazon S3 content. You can also create an origin access identity and add it to your distribution when you create the distribution.

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


### List origin access identity with comments

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
