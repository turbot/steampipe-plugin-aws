select
  account_id,
  allow_users_to_change_password,
  expire_passwords,
  minimum_password_length,
  partition,
  region,
  require_lowercase_characters,
  require_numbers,
  require_symbols,
  require_uppercase_characters
from aws.aws_iam_account_password_policy

