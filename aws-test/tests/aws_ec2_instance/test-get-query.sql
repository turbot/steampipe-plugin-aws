select instance_id, arn, instance_type, monitoring_state, cpu_options_core_count, cpu_options_threads_per_core, ebs_optimized, hypervisor, image_id, tags_src, akas
from aws.aws_ec2_instance
where instance_id = '{{ output.resource_id.value }}'
