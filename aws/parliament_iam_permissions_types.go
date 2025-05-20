package aws

type ParliamentCondition struct {
	Condition   string
	Description string
	Type        string
}

type ParliamentResourceType struct {
	ConditionKeys    []string
	DependentActions []string
	ResourceType     string
}

type ParliamentPrivilege struct {
	AccessLevel   string
	Description   string
	Privilege     string
	ResourceTypes []ParliamentResourceType
}

type ParliamentResource struct {
	Arn           string
	ConditionKeys []string
	Resource      string
}

type ParliamentService struct {
	Conditions  []ParliamentCondition
	Prefix      string
	Privileges  []ParliamentPrivilege
	Resources   []ParliamentResource
	ServiceName string
}

type ParliamentPermissions []ParliamentService
