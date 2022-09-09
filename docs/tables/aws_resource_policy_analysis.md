# Table: aws_resource_policy_analysis

AWS Resource Policy Analysis table returns the access analysis of the IAM policies of the specific resource.

**Important Notes:**

You **_must_** specify a single `policy` and `account_id` in a `where` or `join` clause in order to use this table.

This table answers the following questions:

- Which principals have active grants?
- Which statements grant public, shared or private access?
- What are the [access levels](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_understand-policy-summary-access-level-summaries.html#access_policies_access-level) set at the public, shared and private levels?
- Is the overall visibility of the policy at a public, shared and private level?
- How many accounts, identity providers, services or organizations are referred by the policy?

The analysis table will return overall access granted and it is subdivided into three categories: ["Public", "Shared", "Private"].

Public access

- A policy is public when the policy grants access to a wide range of potentially untrusted accounts due to Pricipal being set to '\*'.
- A policy is public when the policy grants an AWS Service access but has not used any conditions (aws:SourceAccount, aws:SourceArn, aws:SourceOwner) to restrict access to the service.
- A policy is public when an the principal is set to an identity provider and that provider doesn't restrict access by using audience conditions, for example SAML:aud.

Shared access

- A policy is shared if a principal references an account that is differs from the home account where the policy is deploy.
- A policy is shared if a principal is a service with conditions in place.
- A policy is shared if a principal is an identity provider with audience conditions set.

Private access

- A policy is private if the principals in the policy only refer to the home account where the policy is deploy.

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

The evaluator returns an overall report of the access granted by a policy and its statements and not a statement by statement analysis.
The evaluator doesn't sort Principals by the

## Concepts:

**Notes**: For illustration purposes, assuming `111122225555` as the owner account id.

### Access Level

Three levels (i.e. public, shared, private) of access are defined per the evaluation in aws_resource_policy_analysis table.

- `public`: if the policy has at least one `Allow` statement that grants one or more permission to the `*` principal, e.g.,

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
        "Action": ["s3:PutObject", "s3:PutObjectAcl"],
        "Resource": "arn:aws:s3:::EXAMPLE-BUCKET/*"
      }
    ]
  }
  ```

  ```json
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "AllowPublicAccess2",
        "Effect": "Allow",
        "Principal": "*",
        "Action": ["s3:PutObject", "s3:PutObjectAcl"],
        "Resource": "arn:aws:s3:::EXAMPLE-BUCKET/*"
      }
    ]
  }
  ```

- `shared`: If the policy has at least one `Allow` statement that grants one or more permission to the principals (i.e. AWS identities, AWS services) outside of the policy owner account and is not public (as per the above definition), e.g.,

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

  ```json
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "AllowSharedAccess2",
        "Effect": "Allow",
        "Principal": "*",
        "Action": "s3:ListBucket",
        "Resource": "arn:aws:s3:::test-anonymous-access",
        "Condition": {
          "StringEquals": {
            "aws:PrincipalAccount": ["111122223333", "111122224444"]
          }
        }
      }
    ]
  }
  ```

- `private`: If the policy doesn't have any `Allow` statement that grants one or more permission to the principals (i.e. AWS identities, AWS services) outside of the policy owner account, e.g.,

  ```json
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "AllowPrivateAccess1",
        "Effect": "Allow",
        "Principal": {
          "AWS": ["arn:aws:iam::111122225555:root"]
        },
        "Action": "s3:ListBucket",
        "Resource": "arn:aws:s3:::test-anonymous-access"
      }
    ]
  }
  ```

  ```json
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "AllowPrivateAccess2",
        "Effect": "Allow",
        "Principal": "*",
        "Action": "s3:ListBucket",
        "Resource": "arn:aws:s3:::test-anonymous-access",
        "Condition": {
          "StringEquals": {
            "aws:PrincipalAccount": ["111122225555"]
          }
        }
      }
    ]
  }
  ```

### Evaluation constraints

When evaluating statements for public access, the following [condition keys](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_condition-keys.html) are checked:

- `aws:PrincipalAccount`
- `aws:PrincipalArn`
- `aws:PrincipalOrgID`
- `aws:SourceAccount`
- `aws:SourceArn`
- `aws:SourceOwner`

And the following [condition operators](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html) are checked:

- [String condition operators](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_String)
- [Amazon Resource Name (ARN) condition operators](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_ARN)

For each statement, if there are any condition keys with any of the condition operators present then the statement is not considered to be granting public access. An extra check is performed for the `ArnLike` and `StringLike` operators to ensure that the condition key values do not contain `*`.

The inverse condition operators, like `StringNotEquals` and `ArnNotLike`, are not currently evaluated.

### Evaluation logic:

#### Important elements of a policy to deal

- [Principal](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_principal.html)
- [Condition](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition.html)

When a `principal` makes a request to AWS, AWS gathers the request information into a request context. The information is used to evaluate and authorize the request. You can use the `Condition` element of a JSON policy to test specific conditions against the request context.

When a request is submitted, AWS evaluates each condition key in the policy and returns a value of true, false, not present, and occasionally null (an empty data string). A key that is not present in the request is considered a mismatch.

#### [aws:PrincipalAccount](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_condition-keys.html#condition-keys-principalaccount) with `String` operator

The below policy grants public access to sample-bucket objects.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicAccess1",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::sample-bucket/*"
    }
  ]
}
```

```sql
select
  is_public
from
  aws_resource_policy_analysis
where policy = '{"Version":"2012-10-17","Statement":[{"Sid":"PublicAccess1","Effect":"Allow","Principal":"*","Action":"s3:GetObject","Resource":"arn:aws:s3:::sample-bucket/*"}]}' and
account_id = '111122225555';

+-----------+
| is_public |
+-----------+
| true      |
+-----------+
```

But the above statement `PublicAccess1` can be limited to provide access to set of limited aws accounts either by modifying `"Principal": "*"` to limited set of aws accounts like `"Principal": { "AWS": ["arn:aws:iam::111122221111:root", "arn:aws:iam::111122223333:root"] }` or adding the `Condition` element in the policy and using `aws:PrincipalAccount` global condition key to limit to a set of accounts.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicAccess2",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:ListBucket",
      "Resource": "arn:aws:s3:::test-anonymous-access",
      "Condition": {
        "StringEquals": {
          "aws:PrincipalAccount": ["111122221111", "111122223333"]
        }
      }
    }
  ]
}
```

```sql
select
  is_public,
  access_level,
  allowed_principal_account_ids
from
  aws_resource_policy_analysis
where policy = '{"Version":"2012-10-17","Statement":[{"Sid":"PublicAccess2","Effect":"Allow","Principal":"*","Action":"s3:ListBucket","Resource":"arn:aws:s3:::test-anonymous-access","Condition":{"StringEquals":{"aws:PrincipalAccount":["111122221111","111122223333"]}}}]}' and
account_id = '111122225555';

+-----------+--------------+---------------------------------+
| is_public | access_level | allowed_principal_account_ids   |
+-----------+--------------+---------------------------------+
| false     | shared       | ["111122221111","111122223333"] |
+-----------+--------------+---------------------------------+
```

s3://test-anonymous-access/AWSLogs/
aws s3 cp s3://test-anonymous-access/AWSLogs AWSLogs

But

## Refrences

- [The request context](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition.html#AccessPolicyLanguage_RequestContext)
