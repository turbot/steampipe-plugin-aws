select reserved_instance_id, arn, instance_type, instance_state, currency_code, CAST(fixed_price AS varchar), offering_class, scope, CAST(usage_price AS varchar)
from aws.aws_ec2_reserved_instance
where reserved_instance_id = '{{ resourceName }}';
