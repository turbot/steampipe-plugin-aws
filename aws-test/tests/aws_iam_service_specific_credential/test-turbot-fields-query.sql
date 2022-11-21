select title
from aws.aws_iam_service_specific_credential
where service_specific_credential_id = '{{ output.service_specific_credential_id.value }}';
