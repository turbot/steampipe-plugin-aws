select arn, id, name, type, partition from aws.aws_auditmanager_framework
where arn = '{{ output.arn.value }}';
