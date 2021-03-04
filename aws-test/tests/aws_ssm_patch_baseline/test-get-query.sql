select baseline_id, name, description, akas, tags_src, tags, partition, region from aws.aws_ssm_patch_baseline where baseline_id='{{output.resource_id.value}}'
