select title, akas, tags, region, account_id
from aws.aws_ssm_patch_baseline
where baseline_id = '{{ output.resourceId.value }}1a3';