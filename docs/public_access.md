### Docs

## Situation evaluated

1. If `PrincipalArn` contains public arn **`(i.e. if wildcard "*" in the account placeholder in arn)`**.
2. If `SourceArn` contains public arn **`(i.e. if wildcard "*" in the account placeholder in arn)`**.
3. If **`wildcard "*"`** in the statement principals but limited with `PrincipalAccount`, `SourceAccount`, `SourceOwner`, `PrincipalOrgID` and `PrincipalOrgPaths`.
4. Working for multiple condition keys `{"StringEquals":{"aws:PrincipalAccount":"999988887777"},"ArnLike":{"aws:SourceArn":"arn:aws:cloudwatch:us-east-1:*:alarm:*"}}`

   - `{"aws:SourceArn":"arn:aws:cloudwatch:us-east-1:*:alarm:*"}` is a public condition but what makes it `shared` but not `public` is the limitation imposed through `{"StringEquals":{"aws:PrincipalAccount":"999988887777"}` condition.

5. Handles `NotPrincipal With Allow`
6. For `effect = "Deny"`, just not evaluating the statement and marking the statement as blocking public access.
   - When the effect is `Deny` - it doesn't grant access to anyone explictely but only restricts a set of principals from getting access. So, if a policy

### Sample policies

```json
[
  {
    "policy_std": {
      "Statement": [
        {
          "Action": ["s3:getbucketlocation", "s3:listbucket"],
          "Condition": {
            "StringEquals": {
              "aws:principalaccount": ["123456789012"],
              "aws:principalarn": ["arn:aws:iam::111122223333:root"]
            }
          },
          "Effect": "Allow",
          "Principal": {
            "AWS": ["*"]
          },
          "Resource": ["arn:aws:s3:::test"],
          "Sid": "Example permissions"
        },
        {
          "Action": ["s3:getbucketlocation", "s3:listbucket"],
          "Effect": "Deny",
          "Principal": {
            "AWS": [
              "arn:aws:iam::111122223333:root",
              "arn:aws:iam::123456789012:root"
            ]
          },
          "Resource": ["arn:aws:s3:::test"]
        }
      ],
      "Version": "2012-10-17"
    }
  }
]
```

```json
[
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": "dynamodb:GetItem",
        "Resource": "arn:aws:dynamodb:*:*:table/Thread",
        "Condition": {
          "ForAllValues:StringEquals": {
            "dynamodb:Attributes": ["ID", "Message", "Tags"]
          }
        }
      }
    ]
  },
  {
    "Version": "2012-10-17",
    "Statement": {
      "Effect": "Deny",
      "Action": "dynamodb:PutItem",
      "Resource": "arn:aws:dynamodb:*:*:table/Thread",
      "Condition": {
        "ForAnyValue:StringEquals": {
          "dynamodb:Attributes": ["ID", "PostDateTime"]
        }
      }
    }
  },
  {
    "ForAnyValue:StringLike": {
      "aws:PrincipalOrgPaths": [
        "o-a1b2c3d4e5/r-ab12/ou-ab12-33333333/*",
        "o-a1b2c3d4e5/r-ab12/ou-ab12-22222222/*"
      ]
    }
  }
]
```

```json
[
  {
    "ForAnyValue:StringEquals": {
      "aws:PrincipalOrgPaths": [
        "o-a1b2c3d4e5/r-ab12/ou-ab12-11111111/ou-ab12-22222222/"
      ]
    }
  },
  {
    "ForAnyValue:StringLike": {
      "aws:PrincipalOrgPaths": [
        "o-a1b2c3d4e5/r-ab12/ou-ab12-11111111/ou-ab12-22222222/*"
      ]
    }
  },
  {
    "ForAnyValue:StringLike": {
      "aws:PrincipalOrgPaths": [
        "o-a1b2c3d4e5/r-ab12/ou-ab12-11111111/ou-ab12-22222222/*"
      ]
    }
  },
  {
    "ForAnyValue:StringLike": {
      "aws:PrincipalOrgPaths": ["o-a1b2c3d4e5/*"]
    }
  },
  {
    "ForAnyValue:StringLike": {
      "aws:PrincipalOrgPaths": [
        "o-a1b2c3d4e5/r-ab12/ou-ab12-33333333/*",
        "o-a1b2c3d4e5/r-ab12/ou-ab12-22222222/*"
      ]
    }
  }
]
```

### Public policies

```json
{
  "Statement": [
    {
      "Action": "s3:ListBucket",
      "Effect": "Allow",
      "Principal": "*",
      "Resource": "arn:aws:s3:::test-anonymous-access",
      "Sid": "Statement1"
    }
  ],
  "Version": "2012-10-17"
}
```

```json
{
  "Statement": [
    {
      "Action": "s3:ListBucket",
      "Condition": {
        "StringEqualsIfExists": {
          "aws:PrincipalAccount": "111122223333"
        }
      },
      "Effect": "Allow",
      "Principal": "*",
      "Resource": "arn:aws:s3:::test-anonymous-access",
      "Sid": "Statement1"
    }
  ],
  "Version": "2012-10-17"
}
```

```json
{
  "Statement": [
    {
      "Action": "s3:ListBucket",
      "Condition": {
        "StringEquals": {
          "aws:username": "lalit"
        },
        "StringEqualsIfExists": {
          "aws:PrincipalAccount": "111122223333"
        }
      },
      "Effect": "Allow",
      "Principal": "*",
      "Resource": "arn:aws:s3:::test-anonymous-access",
      "Sid": "Statement1"
    }
  ],
  "Version": "2012-10-17"
}
```

```json
{
  "Statement": [
    {
      "Action": "s3:ListBucket",
      "Condition": {
        "ArnLike": {
          "aws:PrincipalArn": "arn:aws:iam::*:user/*\""
        }
      },
      "Effect": "Allow",
      "Principal": "*",
      "Resource": "arn:aws:s3:::test-anonymous-access",
      "Sid": "Statement1"
    }
  ],
  "Version": "2012-10-17"
}
```

```json
{
  "Statement": [
    {
      "Action": "s3:ListBucket",
      "Condition": {
        "ArnLike": {
          "aws:SourceArn": "arn:aws:iam::*:*/*"
        }
      },
      "Effect": "Allow",
      "Principal": "*",
      "Resource": "arn:aws:s3:::test-anonymous-access",
      "Sid": "Statement1"
    }
  ],
  "Version": "2012-10-17"
}
```

```json
{
  "Statement": [
    {
      "Action": "s3:ListBucket",
      "Condition": {
        "ArnLike": {
          "aws:SourceArn": "arn:aws:cloudwatch:us-east-1:*:alarm:*"
        }
      },
      "Effect": "Allow",
      "Principal": "*",
      "Resource": "arn:aws:s3:::test-anonymous-access",
      "Sid": "Statement1"
    }
  ],
  "Version": "2012-10-17"
}
```

```json
{
  "Statement": [
    {
      "Action": "s3:ListBucket",
      "Condition": {
        "ForAnyValue:ArnLike": {
          "aws:PrincipalArn": [
            "arn:aws:iam::*:root",
            "arn:aws:iam::444422223333:root"
          ]
        }
      },
      "Effect": "Allow",
      "Principal": "*",
      "Resource": "arn:aws:s3:::test-anonymous-access",
      "Sid": "Statement1"
    }
  ],
  "Version": "2012-10-17"
}
```

### More Examples

- [Example bucket policies](https://docs.aws.amazon.com/AmazonS3/latest/userguide/example-bucket-policies.html#example-bucket-policies-use-case-2)

- [Example cases for Amazon SNS access control](https://docs.aws.amazon.com/AmazonS3/latest/userguide/access-control-block-public-access.html#access-control-block-public-access-policy-status)
