# Table: aws_resource_policy_analysis

AWS Resource Policy Analysis table analyzes AWS resource or identity policies and returns which Principals have been granted access, the details on the SIDs that provide access to third parties and the type of access given at that access level. It will also calculate the overall access level of the policy which can be used to determine if the policy is too permissive.

The summary of this data will be returned in tabular form.

This table answers the following questions:

- Which principals have active grants?
- Which statements grant public, shared or private access?
- What are the [access levels](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_understand-policy-summary-access-level-summaries.html#access_policies_access-level) set at the public, shared and private levels?
- Is the overall access granted of the policy at a public, shared and private level?
- How many accounts, identity providers, services or organizations are referred by the policy?

The analysis table will return overall access granted and it is subdivided into three categories:0

- Public
- Shared
- Private

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

The table will sort statement SIDs into three different categories, public for statements that grant public access, shared and private for shared and private access respectively.

The table will also give the access level of Functions at each access level. This is can be either none or all of the following values:

- Tagging
- Write
- Read
- List
- Permissions management

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

### List all S3 buckets that have publically accessible resource policies

```sql
select
  r.name,
  r.arn
from
  aws_s3_bucket as r,
  aws_resource_policy_analysis as pa
where
  pa.is_public = true
  and pa.account_id = r.account_id
  and pa.policy = r.policy_std
order by
  r.name
```

### Query the resource policy for all S3 buckets

```sql
select
  r.name,
  pa.access_level,
  pa.allowed_principal_account_ids,
  pa.allowed_principals,
  pa.allowed_principal_services,
  pa.allowed_organization_ids,
  r.arn
from
  aws_s3_bucket as r,
  aws_resource_policy_analysis as pa
where
  pa.account_id = r.account_id
  and pa.policy = r.policy_std
order by
  r.name
```

### Query the assume role policy for all IAM roles

```sql
select
  r.name,
  pa.is_public,
  pa.allowed_principal_account_ids,
  pa.allowed_principals,
  pa.allowed_principal_services,
  pa.allowed_organization_ids,
  pa.allowed_principal_federated_identities,
  r.arn
from
  aws_iam_role as r,
  aws_resource_policy_analysis as pa
where
  pa.account_id = r.account_id
  and pa.policy = r.assume_role_policy_std
order by
  r.name
```

### Get the SIDs that grant public and shared access in all customer managed KMS keys

```sql
select
  right(aliases -> 0 ->> 'AliasName', -6) as alias,
  pa.public_statement_ids,
  pa.shared_statement_ids,
  r.id,
  r.arn
from
  aws_kms_key as r,
  aws_resource_policy_analysis as pa
where
  pa.account_id = r.account_id
  and pa.policy = r.policy_std
  and r.key_manager = 'CUSTOMER'
order by
  r.id
```

### Get the public, shared and private access levels for all Lambda functions

```sql
select
  r.name,
  pa.public_access_levels,
  pa.shared_access_levels,
  pa.private_access_levels,
  r.arn
from
  aws_lambda_function as r,
  aws_resource_policy_analysis as pa
where
  pa.account_id = r.account_id
  and pa.policy = r.policy_std
order by
  r.name
```

### Query all Lambda functions that have tagging and writing capabilities at any access level

```sql
select
  r.name,
  pa.public_access_levels,
  pa.shared_access_levels,
  pa.private_access_levels,
  pa.allowed_principal_account_ids,
  pa.allowed_principals,
  pa.allowed_principal_services,
  pa.allowed_organization_ids,
  r.arn
from
  aws_lambda_function as r,
  aws_resource_policy_analysis as pa
where
  pa.account_id = r.account_id
  and pa.policy = r.policy_std
  and (
    pa.public_access_levels <@ '["Tagging", "Write"]'
    or pa.shared_access_levels <@ '["Tagging", "Write"]'
    or pa.private_access_levels <@ '["Tagging", "Write"]'
  )
order by
  r.name
```

### Return all EFS resources that have shared access and are accessible from accounts outside a trusted accounts list

```sql
select
  pa.allowed_principal_account_ids
from
  aws_efs_file_system as r,
  aws_resource_policy_analysis as pa
where
  pa.account_id = r.account_id
  and pa.policy = r.policy_std
  and not pa.allowed_principal_account_ids <@ '["<Account 1>", "<Account 2>"]'
  and pa.is_public = false
  and jsonb_array_length(pa.shared_statement_ids) > 0
```

### Analyze a resource policy query

```sql
select
  pa.is_public,
  pa.allowed_principal_account_ids,
  pa.allowed_principals,
  pa.allowed_principal_services,
  pa.allowed_organization_ids
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
        "Sid": "AllowedPricipal",
        "Effect": "Allow",
        "Principal": { "AWS": "arn:aws:iam::111122223333:root" },
        "Resource": "arn:aws:s3:::test-bucket",
        "Action": "s3:*"
      },
      {
        "Sid": "AllowedService",
        "Effect": "Allow",
        "Principal": { "Service": "ec2.amazonaws.com" },
        "Resource": "arn:aws:s3:::test-bucket",
        "Action": "s3:*",
        "Condition": { "StringEquals": { "aws:SourceAccount": "555566667777" } }
      },
      {
        "Sid": "AllowedOrganization",
        "Effect": "Allow",
        "Principal": { "AWS": "*" },
        "Action": "s3:*",
        "Resource": "arn:aws:s3:::test-bucket",
        "Condition": { "StringEquals": { "aws:PrincipalOrgID": "o-aaabbbccc123" } }
      }
    ]
  }
  '
```

### Analyze a trust policy query

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
    "Statement": [
      {
        "Sid": "AwsPrincipal",
        "Effect": "Allow",
        "Principal": { "AWS": "arn:aws:iam::444455556666:root" },
        "Action": "sts:AssumeRole"
      },
      {
        "Sid": "Federated",
        "Effect": "Allow",
        "Principal": { "Federated": "arn:aws:iam::111122223333:saml-provider/SSO" },
        "Action": "sts:AssumeRoleWithSAML",
        "Condition": { "StringEquals": { "SAML:aud": "aud" } }
      },
      {
        "Sid": "Service",
        "Effect": "Allow",
        "Principal": { "Service": "ec2.amazonaws.com" },
        "Action": "sts:AssumeRole",
        "Condition": { "StringEquals": { "aws:SourceAccount": "011026983618" } }
      },
      {
        "Sid": "Organization",
        "Effect": "Allow",
        "Principal": { "AWS": "*" },
        "Action": "sts:AssumeRole",
        "Condition": { "StringEquals": { "aws:PrincipalOrgID": "o-aaabbbccc123" } }
      },
      {
        "Sid": "WebIdentity",
        "Effect": "Allow",
        "Principal": { "Federated": "accounts.google.com" },
        "Action": "sts:AssumeRoleWithWebIdentity",
        "Condition": { "StringEquals": { "accounts.google.com:aud": "aud" } }
      }
    ]
  }
  '
```
