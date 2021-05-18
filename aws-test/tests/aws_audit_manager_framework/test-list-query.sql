select arn, id, name, type, partition from aws.aws_audit_manager_framework
where arn = '{{ output.arn.value }}';
