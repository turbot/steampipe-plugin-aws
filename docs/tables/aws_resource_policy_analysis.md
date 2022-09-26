# Table: aws_resource_policy_analysis

AWS Resource Policy Analysis table returns the access analysis of the IAM policies of the specific resource.

This table answers the following questions:

- Which principals have active grants?
- Which statements grant public, shared or private access?
- What are the [access levels](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_understand-policy-summary-access-level-summaries.html#access_policies_access-level) set at the public, shared and private levels?
- Is the overall access granted of the policy at a public, shared and private level?
- How many accounts, identity providers, services or organizations are referred by the policy?

The analysis table will return overall access granted and it is subdivided into three categories: ["Public", "Shared", "Private"].

Public access has one of the following characteristics:

- A policy is public when the policy grants access to a wide range of potentially untrusted accounts due to Pricipal being set to `*`.
- A policy is public when the policy grants an AWS Service access but has not used any conditions (aws:SourceAccount, aws:SourceArn, aws:SourceOwner) to restrict access to the service.
- A policy is public when an the principal is set to an identity provider and that provider doesn't restrict access by using audience conditions, for example SAML:aud.

For example:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowPublicAccess1",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::test-anonymous-access"
    }
  ]
}
```

Shared access has one of the following characteristics:

- A policy is shared if a principal references an account that differs from the home account where the policy is deployed.
- A policy is shared if a principal is a service with conditions in place.
- A policy is shared if a principal is an identity provider with audience conditions set.

  ```json
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "AllowSharedAccess1",
        "Effect": "Allow",
        "Principal": {
          "AWS": [
            "arn:aws:iam::111122223333:root",
            "arn:aws:iam::111122224444:root"
          ]
        },
        "Action": "s3:ListBucket",
        "Resource": "arn:aws:s3:::test-anonymous-access"
      }
    ]
  }
  ```

Private access

- A policy is private if the principals in the policy only refers to the home account where the policy is deployed.

For example, if the home account is `111122221111` the following policy will be returned as private:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowPrivateAccess1",
      "Effect": "Allow",
      "Principal": { "AWS": "arn:aws:iam::111122221111:root" },
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::test-anonymous-access"
    }
  ]
}
```

Principals are organized by the table and returned in their own columns.

The table will sort statement IDs into three different categories, public for statements that grant public access, shared and private for shared and private access respectively.

## Limitations

The table evaulates a subset of conditions at present:

- `aws:PrincipalAccount`
- `aws:PrincipalArn`
- `aws:PrincipalOrgID`
- `aws:SourceAccount`
- `aws:SourceArn`
- `aws:SourceOwner`

And the following [condition operators](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html) are checked:

- [String condition operators](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_String)
- [Amazon Resource Name (ARN) condition operators](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_ARN)

The inverse condition operators, like `StringNotEquals` and `ArnNotLike`, are not currently evaluated.
If a condition operator ends if `IfEquals` the table will ignre this condition in its evaulation.

**Important Notes:**

You **_must_** specify a single `policy` and `account_id` in a `where` or `join` clause in order to use this table.

## Examples

### Basic info

```sql
select
  pa.is_public,
  pa.allowed_principal_account_ids,
  pa.allowed_principals,
  pa.allowed_principal_services,
  pa.allowed_organization_ids,
  pa.allowed_principal_federated_identities
from
  aws_resource_policy_analysis as pa
where
  account_id = '111122223333'
  and policy = '
{
  "Version": "2012-10-17",
  "Id": "Policy1658140668960",
  "Statement": [
    {
      "Sid": "AllowedService",
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::omero-resource-policy-bucket",
      "Condition": { "StringEquals": { "aws:SourceAccount": "011026983618" } }
    },
    {
      "Sid": "AllowedPricipal",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::034519905315:root"
      },
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::omero-resource-policy-bucket"
    },
    {
      "Sid": "AllowedOrganization",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::omero-resource-policy-bucket",
      "Condition": { "StringEquals": { "aws:PrincipalOrgID": "o-c3a5y4wd52" } }
    }
  ]
}
  '
```

### Query to analyse policy content for public grants

```sql
select
  *
from
  aws_resource_policy_analysis as ps
where
  account_id = '111122223333'
  and policy = '
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "s3:DeleteBucket",
          "Principal": "*",
          "Resource": "arn:aws:s3:::test-bucket"
        }
      ]
    }
  '
```

### Query the resource policy for an S3 bucket

```sql
select
  pa.access_level,
  pa.
  r.region,
  r.account_id
from
  aws_s3_bucket as r,
  aws_resource_policy_analysis as pa
where
  pa.account_id = r.account_id
  and pa.policy = r.policy_std
group by
  resource,
  title,
  r.region,
  r.account_id
```
