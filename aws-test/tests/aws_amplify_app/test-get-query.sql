select
    arn,
    app_id,
    tags,
    environment_variables,
    enable_branch_auto_deletion,
    enable_branch_auto_build,
    enable_basic_auth,
    description,
    build_spec,
    custom_rules
from aws.aws_amplify_app
where app_id = '{{ output.id.value }}';