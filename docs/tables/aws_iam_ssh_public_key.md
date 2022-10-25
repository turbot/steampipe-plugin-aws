# Table: aws_iam_ssh_public_key

The SSH public keys returned by this table are used for authenticating the IAM user to an CodeCommit repository.

## Examples

### List of SSH public keys with their corresponding user name and date of creation

```sql
select
  ssh_public_key_id,
  ssh_public_key_body_pem,
  ssh_public_key_body_rsa,
  user_name,
  upload_date
from
  aws_iam_ssh_public_key;
```


### List of SSH public keys which are inactive

```sql
select
  ssh_public_key_id,
  user_name,
  status
from
  aws_iam_ssh_public_key
where
  status = 'Inactive';
```


### Access key count by user name

```sql
select
  user_name,
  count (ssh_public_key_id) as ssh_public_key_count
from
  aws_iam_ssh_public_key
group by
  user_name;
```
