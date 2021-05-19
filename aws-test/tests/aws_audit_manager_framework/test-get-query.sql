select arn, id, name, type, partition from aws.aws_audit_manager_framework
where id = '{{ output.id.value }}' and type = '{{ output.type.value }}' and region = '{{ output.aws_region.value }}';