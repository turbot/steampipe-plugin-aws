import subprocess


def escape_string(input_string):
    return input_string.replace("\"", "\\\"")


def write_conditions(conditions, go_file):
    if len(conditions) == 0:
        go_file.write("""Conditions: []ParliamentCondition{},
""")
    else:
        go_file.write("""Conditions: []ParliamentCondition{
""")
        for condition in conditions:
            go_file.write("""{
""")
            go_file.write("""Condition: \"{0}\",
""".format(escape_string(condition["condition"])))
            go_file.write("""Description: \"{0}\",
""".format(escape_string(condition["description"])))
            go_file.write("""Type: \"{0}\",
""".format(escape_string(condition["type"])))
            go_file.write("""},
""")
        go_file.write("""},
""")


def write_condition_keys(condition_keys, go_file):
    if len(condition_keys) == 0:
        go_file.write("""ConditionKeys: []string{},
""")
    else:
        go_file.write("""ConditionKeys: []string{
""")
        for condition_key in condition_keys:
            go_file.write("""\"{0}\",
""".format(escape_string(condition_key)))
        go_file.write("""},
""")


def write_resource_type_dependent_actions(dependent_actions, go_file):
    if len(dependent_actions) == 0:
        go_file.write("""DependentActions: []string{},
""")
    else:
        go_file.write("""DependentActions: []string{
""")
        for dependent_action in dependent_actions:
            go_file.write("""\"{0}\",
""".format(escape_string(dependent_action)))
        go_file.write("""},
""")


def write_resource_types(resource_types, go_file):
    if len(resource_types) == 0:
        go_file.write("""ResourceTypes: []ParliamentResourceType{},
""")
    else:
        go_file.write("""ResourceTypes: []ParliamentResourceType{
""")
        for resource_type in resource_types:
            go_file.write("""{
""")
            write_condition_keys(resource_type.get("condition_keys", []), go_file)
            write_resource_type_dependent_actions(resource_type.get("dependent_actions", []), go_file)
            go_file.write("""ResourceType: \"{0}\",
""".format(escape_string(resource_type["resource_type"])))
            go_file.write("""},
""")
        go_file.write("""},
""")


def write_privileges(privileges, go_file):
    if len(privileges) == 0:
        go_file.write("""Privileges: []ParliamentPrivilege{},
""")
    else:
        go_file.write("""Privileges: []ParliamentPrivilege{
""")
        for privilege in privileges:
            go_file.write("""{
""")
            go_file.write("""AccessLevel: \"{0}\",
""".format(escape_string(privilege["access_level"])))
            go_file.write("""Description: \"{0}\",
""".format(escape_string(privilege["description"])))
            go_file.write("""Privilege: \"{0}\",
""".format(escape_string(privilege["privilege"])))
            write_resource_types(privilege.get("resource_types", []), go_file)
            go_file.write("""},
""")
        go_file.write("""},
""")


def write_resources(resources, go_file):
    if len(resources) == 0:
        go_file.write("""Resources: []ParliamentResource{},
""")
    else:
        go_file.write("""Resources: []ParliamentResource{
""")
        for resource in resources:
            go_file.write("""{
""")
            go_file.write("""Arn: \"{0}\",
""".format(escape_string(resource["arn"])))
            write_condition_keys(resource.get("condition_keys", []), go_file)
            go_file.write("""Resource: \"{0}\",
""".format(escape_string(resource["resource"])))
            go_file.write("""},
""")
        go_file.write("""},
""")


def write_service(service, go_file):
    go_file.write("""ParliamentService{
""")
    write_conditions(service.get("conditions", []), go_file)
    go_file.write("""Prefix: \"{0}\",
""".format(escape_string(service["prefix"])))
    write_privileges(service.get("privileges", []), go_file)
    write_resources(service.get("resources", []), go_file)
    go_file.write("""ServiceName: \"{0}\",
    """.format(escape_string(service["service_name"])))
    go_file.write("""},
""")


def generate_go_file(iam_permissions):
    with open('../../aws/parliament_iam_permissions.go', 'w') as go_file:
        go_file.write("""//go:build !dev

package aws

func getParliamentIamPermissions() ParliamentPermissions {
permissions := ParliamentPermissions{
""")
        for service in iam_permissions:
            write_service(service, go_file)
        go_file.write("""}
return permissions
}
""")
    subprocess.run(["gofmt", "-w", "../../aws/parliament_iam_permissions.go"])
