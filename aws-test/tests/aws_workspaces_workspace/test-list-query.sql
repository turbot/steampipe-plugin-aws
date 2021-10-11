select workspace_id, state
from aws.aws_workspaces_workspace
where akas::text = '["{{ output.resource_aka.value }}"]';