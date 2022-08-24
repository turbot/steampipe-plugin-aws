select *
from aws_dms_replication_instance
where arn = '{{ output.resource_aka.value }}1p000';