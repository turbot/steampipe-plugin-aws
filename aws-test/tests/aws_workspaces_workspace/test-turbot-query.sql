select akas
from aws.aws_workspaces_workspace
where workspace_id = '{{ output.id.value }}';