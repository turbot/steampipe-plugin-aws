select title, tags, akas from aws.aws_amplify_app
where app_id = '{{ output.id.value }}';