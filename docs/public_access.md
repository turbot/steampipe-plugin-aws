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
