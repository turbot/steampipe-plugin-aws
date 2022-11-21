select name, akas
from aws_sagemaker_app
where user_profile_name = '{{ output.user_profile_name.value }}';
