select title, akas, tags, region, account_id
from aws.aws_sagemaker_app
where name = 'dummy-{{ resourceName }}'
and user_profile_name = '{{ output.user_profile_name.value }}';