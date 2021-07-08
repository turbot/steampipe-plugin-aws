select arn, auth_type, server_type
from aws.aws_codebuild_source_credential
where arn = '{{ output.arn.value }}';