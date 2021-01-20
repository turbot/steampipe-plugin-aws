# Table: aws_iam_access_key

Access keys are long-term credentials for an IAM user or the AWS account root user. You can use access keys to sign programmatic requests to the AWS CLI or AWS API (directly or using the AWS SDK).

## Examples

### List of access keys with their corresponding user name and date of creation

```sql
select
  access_key_id,
  user_name,
  create_date
from
  aws_iam_access_key;
```


### List of access keys which are inactive

```sql
select
  access_key_id,
  user_name,
  status
from
  aws_iam_access_key
where
  status = 'Inactive';
```


### Access key count by user name

```sql
select
  user_name,
  count (access_key_id) as access_key_count
from
  aws_iam_access_key
group by
  user_name;
```