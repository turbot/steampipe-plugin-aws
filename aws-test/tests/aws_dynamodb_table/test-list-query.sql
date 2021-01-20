select table_arn, name, attribute_definitions, write_capacity, read_capacity, key_schema
from aws.aws_dynamodb_table
where akas::text = '["{{output.resource_aka.value}}"]'
