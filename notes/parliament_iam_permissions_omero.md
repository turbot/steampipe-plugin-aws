# Code

```go
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

func getParliamentIamPermissions() ParliamentPermissions {
	permissions := ParliamentPermissions{
		ParliamentService{
			Conditions: []ParliamentCondition{
				{
					Condition:   "aws:RequestTag/${TagKey}",
					Description: "Filters access by a tag key and value pair that is allowed in the request",
					Type:        "String",
				},
				{
					Condition:   "aws:ResourceTag/",
					Description: "Filters access by the preface string for a tag key and value pair that are attached to a resource",
					Type:        "String",
				},
				{
					Condition:   "aws:ResourceTag/${TagKey}",
					Description: "Filters access by a tag key and value pair of a resource",
					Type:        "String",
				},
				{
					Condition:   "aws:TagKeys",
					Description: "Filters access by a list of tag keys that are allowed in the request",
					Type:        "String",
				},
				{
					Condition:   "ec2:AccepterVpc",
					Description: "Filters access by the ARN of an accepter VPC in a VPC peering connection",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:AssociatePublicIpAddress",
					Description: "Filters access by whether the user wants to associate a public IP address with the instance",
					Type:        "Bool",
				},
				{
					Condition:   "ec2:Attribute/${AttributeName}",
					Description: "Filters access by an attribute being set on a resource",
					Type:        "String",
				},
				{
					Condition:   "ec2:AuthenticationType",
					Description: "Filters access by the authentication type for the VPN tunnel endpoints",
					Type:        "String",
				},
				{
					Condition:   "ec2:AuthorizedService",
					Description: "Filters access by the AWS service that has permission to use a resource",
					Type:        "String",
				},
				{
					Condition:   "ec2:AuthorizedUser",
					Description: "Filters access by an IAM principal that has permission to use a resource",
					Type:        "String",
				},
				{
					Condition:   "ec2:AutoPlacement",
					Description: "Filters access by the Auto Placement properties of a Dedicated Host",
					Type:        "String",
				},
				{
					Condition:   "ec2:AvailabilityZone",
					Description: "Filters access by the name of an Availability Zone in an AWS Region",
					Type:        "String",
				},
				{
					Condition:   "ec2:CapacityReservationFleet",
					Description: "Filters access by the ARN of the Capacity Reservation Fleet",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:ClientRootCertificateChainArn",
					Description: "Filters access by the ARN of the client root certificate chain",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:CloudwatchLogGroupArn",
					Description: "Filters access by the ARN of the CloudWatch Logs log group",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:CloudwatchLogStreamArn",
					Description: "Filters access by the ARN of the CloudWatch Logs log stream",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:CreateAction",
					Description: "Filters access by the name of a resource-creating API action",
					Type:        "String",
				},
				{
					Condition:   "ec2:DPDTimeoutSeconds",
					Description: "Filters access by the duration after which DPD timeout occurs on a VPN tunnel",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:DirectoryArn",
					Description: "Filters access by the ARN of the directory",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:EbsOptimized",
					Description: "Filters access by whether the instance is enabled for EBS optimization",
					Type:        "Bool",
				},
				{
					Condition:   "ec2:ElasticGpuType",
					Description: "Filters access by the type of Elastic Graphics accelerator",
					Type:        "String",
				},
				{
					Condition:   "ec2:Encrypted",
					Description: "Filters access by whether the EBS volume is encrypted",
					Type:        "Bool",
				},
				{
					Condition:   "ec2:GatewayType",
					Description: "Filters access by the gateway type for a VPN endpoint on the AWS side of a VPN connection",
					Type:        "String",
				},
				{
					Condition:   "ec2:HostRecovery",
					Description: "Filters access by whether host recovery is enabled for a Dedicated Host",
					Type:        "String",
				},
				{
					Condition:   "ec2:IKEVersions",
					Description: "Filters access by the internet key exchange (IKE) versions that are permitted for a VPN tunnel",
					Type:        "String",
				},
				{
					Condition:   "ec2:ImageType",
					Description: "Filters access by the type of image (machine, aki, or ari)",
					Type:        "String",
				},
				{
					Condition:   "ec2:InsideTunnelCidr",
					Description: "Filters access by the range of inside IP addresses for a VPN tunnel",
					Type:        "String",
				},
				{
					Condition:   "ec2:InstanceMarketType",
					Description: "Filters access by the market or purchasing option of an instance (on-demand or spot)",
					Type:        "String",
				},
				{
					Condition:   "ec2:InstanceProfile",
					Description: "Filters access by the ARN of an instance profile",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:InstanceType",
					Description: "Filters access by the type of instance",
					Type:        "String",
				},
				{
					Condition:   "ec2:IsLaunchTemplateResource",
					Description: "Filters access by whether users are able to override resources that are specified in the launch template",
					Type:        "Bool",
				},
				{
					Condition:   "ec2:KeyPairName",
					Description: "Filters access by key pair name",
					Type:        "String",
				},
				{
					Condition:   "ec2:LaunchTemplate",
					Description: "Filters access by the ARN of a launch template",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:MetadataHttpEndpoint",
					Description: "Filters access by whether the HTTP endpoint is enabled for the instance metadata service",
					Type:        "String",
				},
				{
					Condition:   "ec2:MetadataHttpPutResponseHopLimit",
					Description: "Filters access by the allowed number of hops when calling the instance metadata service",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:MetadataHttpTokens",
					Description: "Filters access by whether tokens are required when calling the instance metadata service (optional or required)",
					Type:        "String",
				},
				{
					Condition:   "ec2:NewInstanceProfile",
					Description: "Filters access by the ARN of the instance profile being attached",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:OutpostArn",
					Description: "Filters access by the ARN of the Outpost",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:Owner",
					Description: "Filters access by the owner of the resource (amazon, aws-marketplace, or an AWS account ID)",
					Type:        "String",
				},
				{
					Condition:   "ec2:ParentSnapshot",
					Description: "Filters access by the ARN of the parent snapshot",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:ParentVolume",
					Description: "Filters access by the ARN of the parent volume from which the snapshot was created",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:Permission",
					Description: "Filters access by the type of permission for a resource (INSTANCE-ATTACH or EIP-ASSOCIATE)",
					Type:        "String",
				},
				{
					Condition:   "ec2:Phase1DHGroupNumbers",
					Description: "Filters access by the Diffie-Hellman group numbers that are permitted for a VPN tunnel for the phase 1 IKE negotiations",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:Phase1EncryptionAlgorithms",
					Description: "Filters access by the encryption algorithms that are permitted for a VPN tunnel for the phase 1 IKE negotiations",
					Type:        "String",
				},
				{
					Condition:   "ec2:Phase1IntegrityAlgorithms",
					Description: "Filters access by the integrity algorithms that are permitted for a VPN tunnel for the phase 1 IKE negotiations",
					Type:        "String",
				},
				{
					Condition:   "ec2:Phase1LifetimeSeconds",
					Description: "Filters access by the lifetime in seconds for phase 1 of the IKE negotiations for a VPN tunnel",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:Phase2DHGroupNumbers",
					Description: "Filters access by the Diffie-Hellman group numbers that are permitted for a VPN tunnel for the phase 2 IKE negotiations",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:Phase2EncryptionAlgorithms",
					Description: "Filters access by the encryption algorithms that are permitted for a VPN tunnel for the phase 2 IKE negotiations",
					Type:        "String",
				},
				{
					Condition:   "ec2:Phase2IntegrityAlgorithms",
					Description: "Filters access by the integrity algorithms that are permitted for a VPN tunnel for the phase 2 IKE negotiations",
					Type:        "String",
				},
				{
					Condition:   "ec2:Phase2LifetimeSeconds",
					Description: "Filters access by the lifetime in seconds for phase 2 of the IKE negotiations for a VPN tunnel",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:PlacementGroup",
					Description: "Filters access by the ARN of the placement group",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:PlacementGroupStrategy",
					Description: "Filters access by the instance placement strategy used by the placement group (cluster, spread, or partition)",
					Type:        "String",
				},
				{
					Condition:   "ec2:PresharedKeys",
					Description: "Filters access by the pre-shared key (PSK) used to establish the initial IKE security association between a virtual private gateway and a customer gateway",
					Type:        "String",
				},
				{
					Condition:   "ec2:ProductCode",
					Description: "Filters access by the product code that is associated with the AMI",
					Type:        "String",
				},
				{
					Condition:   "ec2:Public",
					Description: "Filters access by whether the image has public launch permissions",
					Type:        "Bool",
				},
				{
					Condition:   "ec2:Quantity",
					Description: "Filters access by the number of Dedicated Hosts in a request",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:Region",
					Description: "Filters access by the name of the AWS Region",
					Type:        "String",
				},
				{
					Condition:   "ec2:RekeyFuzzPercentage",
					Description: "Filters access by the percentage of increase of the rekey window (determined by the rekey margin time) within which the rekey time is randomly selected for a VPN tunnel",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:RekeyMarginTimeSeconds",
					Description: "Filters access by the margin time before the phase 2 lifetime expires for a VPN tunnel",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:RequesterVpc",
					Description: "Filters access by the ARN of a requester VPC in a VPC peering connection",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:ReservedInstancesOfferingType",
					Description: "Filters access by the payment option of the Reserved Instance offering (No Upfront, Partial Upfront, or All Upfront)",
					Type:        "String",
				},
				{
					Condition:   "ec2:ResourceTag/",
					Description: "Filters access by the preface string for a tag key and value pair that are attached to a resource",
					Type:        "String",
				},
				{
					Condition:   "ec2:ResourceTag/${TagKey}",
					Description: "Filters access by a tag key and value pair of a resource",
					Type:        "String",
				},
				{
					Condition:   "ec2:RoleDelivery",
					Description: "Filters access by the version of the instance metadata service for retrieving IAM role credentials for EC2",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:RootDeviceType",
					Description: "Filters access by the root device type of the instance (ebs or instance-store)",
					Type:        "String",
				},
				{
					Condition:   "ec2:RoutingType",
					Description: "Filters access by the routing type for the VPN connection",
					Type:        "String",
				},
				{
					Condition:   "ec2:SamlProviderArn",
					Description: "Filters access by the ARN of the IAM SAML identity provider",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:ServerCertificateArn",
					Description: "Filters access by the ARN of the server certificate",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:SnapshotTime",
					Description: "Filters access by the initiation time of a snapshot",
					Type:        "String",
				},
				{
					Condition:   "ec2:SourceInstanceARN",
					Description: "Filters access by the ARN of the instance from which the request originated",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:SourceOutpostArn",
					Description: "Filters access by the ARN of the Outpost from which the request originated",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:Subnet",
					Description: "Filters access by the ARN of the subnet",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:Tenancy",
					Description: "Filters access by the tenancy of the VPC or instance (default, dedicated, or host)",
					Type:        "String",
				},
				{
					Condition:   "ec2:VolumeIops",
					Description: "Filters access by the the number of input/output operations per second (IOPS) provisioned for the volume",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:VolumeSize",
					Description: "Filters access by the size of the volume, in GiB",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:VolumeThroughput",
					Description: "Filters access by the throughput of the volume, in MiBps",
					Type:        "Numeric",
				},
				{
					Condition:   "ec2:VolumeType",
					Description: "Filters access by the type of volume (gp2, gp3, io1, io2, st1, sc1, or standard)",
					Type:        "String",
				},
				{
					Condition:   "ec2:Vpc",
					Description: "Filters access by the ARN of the VPC",
					Type:        "ARN",
				},
				{
					Condition:   "ec2:VpceServiceName",
					Description: "Filters access by the name of the VPC endpoint service",
					Type:        "String",
				},
				{
					Condition:   "ec2:VpceServiceOwner",
					Description: "Filters access by the service owner of the VPC endpoint service (amazon, aws-marketplace, or an AWS account ID)",
					Type:        "String",
				},
				{
					Condition:   "ec2:VpceServicePrivateDnsName",
					Description: "Filters access by the private DNS name of the VPC endpoint service",
					Type:        "String",
				},
			},
			Prefix: "ec2",
			Privileges: []ParliamentPrivilege{
				{
					AccessLevel: "Write",
					Description: "Grants permission to accept a Convertible Reserved Instance exchange quote",
					Privilege:   "AcceptReservedInstancesExchangeQuote",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:InstanceType",
								"ec2:Region",
								"ec2:ReservedInstancesOfferingType",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "reserved-instances*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to accept a request to associate subnets with a transit gateway multicast domain",
					Privilege:   "AcceptTransitGatewayMulticastDomainAssociations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to accept a transit gateway peering attachment request",
					Privilege:   "AcceptTransitGatewayPeeringAttachment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to accept a request to attach a VPC to a transit gateway",
					Privilege:   "AcceptTransitGatewayVpcAttachment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to accept one or more interface VPC endpoint connections to your VPC endpoint service",
					Privilege:   "AcceptVpcEndpointConnections",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to accept a VPC peering connection request",
					Privilege:   "AcceptVpcPeeringConnection",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AccepterVpc",
								"ec2:Region",
								"ec2:RequesterVpc",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-peering-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to advertise an IP address range that is provisioned for use in AWS through bring your own IP addresses (BYOIP)",
					Privilege:   "AdvertiseByoipCidr",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to allocate an Elastic IP address (EIP) to your account",
					Privilege:   "AllocateAddress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:ResourceTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "ipv4pool-ec2",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to allocate a Dedicated Host to your account",
					Privilege:   "AllocateHosts",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to apply a security group to the association between a Client VPN endpoint and a target network",
					Privilege:   "ApplySecurityGroupsToClientVpnTargetNetwork",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to assign one or more IPv6 addresses to a network interface",
					Privilege:   "AssignIpv6Addresses",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to assign one or more secondary private IP addresses to a network interface",
					Privilege:   "AssignPrivateIpAddresses",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate an Elastic IP address (EIP) with an instance or a network interface",
					Privilege:   "AssociateAddress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "elastic-ip",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate a target network with a Client VPN endpoint",
					Privilege:   "AssociateClientVpnTargetNetwork",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate or disassociate a set of DHCP options with a VPC",
					Privilege:   "AssociateDhcpOptions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate an ACM certificate with an IAM role to be used in an EC2 Enclave",
					Privilege:   "AssociateEnclaveCertificateIamRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "certificate*",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate an IAM instance profile with a running or stopped instance",
					Privilege:   "AssociateIamInstanceProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:NewInstanceProfile",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{
								"iam:PassRole",
							},
							ResourceType: "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate one or more targets with an event window",
					Privilege:   "AssociateInstanceEventWindow",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "instance-event-window*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate a subnet or gateway with a route table",
					Privilege:   "AssociateRouteTable",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "internet-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate a CIDR block with a subnet",
					Privilege:   "AssociateSubnetCidrBlock",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate an attachment and list of subnets with a transit gateway multicast domain",
					Privilege:   "AssociateTransitGatewayMulticastDomain",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate an attachment with a transit gateway route table",
					Privilege:   "AssociateTransitGatewayRouteTable",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate a branch network interface with a trunk network interface",
					Privilege:   "AssociateTrunkInterface",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate a CIDR block with a VPC",
					Privilege:   "AssociateVpcCidrBlock",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "ipv6pool-ec2",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to link an EC2-Classic instance to a ClassicLink-enabled VPC through one or more of the VPC's security groups",
					Privilege:   "AttachClassicLinkVpc",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to attach an internet gateway to a VPC",
					Privilege:   "AttachInternetGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "internet-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to attach a network interface to an instance",
					Privilege:   "AttachNetworkInterface",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to attach an EBS volume to a running or stopped instance and expose it to the instance with the specified device name",
					Privilege:   "AttachVolume",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to attach a virtual private gateway to a VPC",
					Privilege:   "AttachVpnGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add an inbound authorization rule to a Client VPN endpoint",
					Privilege:   "AuthorizeClientVpnIngress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add one or more outbound rules to a VPC security group",
					Privilege:   "AuthorizeSecurityGroupEgress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add one or more inbound rules to a security group",
					Privilege:   "AuthorizeSecurityGroupIngress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to bundle an instance store-backed Windows instance",
					Privilege:   "BundleInstance",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to cancel a bundling operation",
					Privilege:   "CancelBundleTask",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to cancel a Capacity Reservation and release the reserved capacity",
					Privilege:   "CancelCapacityReservation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:CapacityReservationFleet",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to cancel an active conversion task",
					Privilege:   "CancelConversionTask",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to cancel an active export task",
					Privilege:   "CancelExportTask",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "export-image-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "export-instance-task",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to cancel an in-process import virtual machine or import snapshot task",
					Privilege:   "CancelImportTask",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "import-image-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "import-snapshot-task",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to cancel a Reserved Instance listing on the Reserved Instance Marketplace",
					Privilege:   "CancelReservedInstancesListing",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to cancel one or more Spot Fleet requests",
					Privilege:   "CancelSpotFleetRequests",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "spot-fleet-request*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to cancel one or more Spot Instance requests",
					Privilege:   "CancelSpotInstanceRequests",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "spot-instances-request*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to determine whether an owned product code is associated with an instance",
					Privilege:   "ConfirmProductInstance",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to copy a source Amazon FPGA image (AFI) to the current Region. Resource-level permissions specified for this action apply to the new AFI only. They do not apply to the source AFI",
					Privilege:   "CopyFpgaImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"ec2:Owner",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "fpga-image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to copy an Amazon Machine Image (AMI) from a source Region to the current Region. Resource-level permissions specified for this action apply to the new AMI only. They do not apply to the source AMI",
					Privilege:   "CopyImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"ec2:Owner",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to copy a point-in-time snapshot of an EBS volume and store it in Amazon S3. Resource-level permissions specified for this action apply to the new snapshot only. They do not apply to the source snapshot",
					Privilege:   "CopySnapshot",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:OutpostArn",
								"ec2:Region",
								"ec2:SourceOutpostArn",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a Capacity Reservation",
					Privilege:   "CreateCapacityReservation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:CapacityReservationFleet",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a carrier gateway and provides CSP connectivity to VPC customers",
					Privilege:   "CreateCarrierGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
								"ec2:Tenancy",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "carrier-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a Client VPN endpoint",
					Privilege:   "CreateClientVpnEndpoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add a network route to a Client VPN endpoint's route table",
					Privilege:   "CreateClientVpnRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a customer gateway, which provides information to AWS about your customer gateway device",
					Privilege:   "CreateCustomerGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "customer-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a default subnet in a specified Availability Zone in a default VPC",
					Privilege:   "CreateDefaultSubnet",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a default VPC with a default subnet in each Availability Zone",
					Privilege:   "CreateDefaultVpc",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a set of DHCP options for a VPC",
					Privilege:   "CreateDhcpOptions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "dhcp-options*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an egress-only internet gateway for a VPC",
					Privilege:   "CreateEgressOnlyInternetGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "egress-only-internet-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to launch an EC2 Fleet",
					Privilege:   "CreateFleet",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "fleet*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create one or more flow logs to capture IP traffic for a network interface",
					Privilege:   "CreateFlowLogs",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{
								"iam:PassRole",
							},
							ResourceType: "vpc-flow-log*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an Amazon FPGA Image (AFI) from a design checkpoint (DCP)",
					Privilege:   "CreateFpgaImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "fpga-image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an Amazon EBS-backed AMI from a stopped or running Amazon EBS-backed instance",
					Privilege:   "CreateImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an event window in which scheduled events for the associated Amazon EC2 instances can run",
					Privilege:   "CreateInstanceEventWindow",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "instance-event-window*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to export a running or stopped instance to an Amazon S3 bucket",
					Privilege:   "CreateInstanceExportTask",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "export-instance-task*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an internet gateway for a VPC",
					Privilege:   "CreateInternetGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "internet-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a 2048-bit RSA key pair",
					Privilege:   "CreateKeyPair",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:KeyPairName",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a launch template",
					Privilege:   "CreateLaunchTemplate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:PlacementGroupStrategy",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "placement-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new version of a launch template",
					Privilege:   "CreateLaunchTemplateVersion",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:PlacementGroupStrategy",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "placement-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a static route for a local gateway route table",
					Privilege:   "CreateLocalGatewayRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-virtual-interface-group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to associate a VPC with a local gateway route table",
					Privilege:   "CreateLocalGatewayRouteTableVpcAssociation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table-vpc-association*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a managed prefix list",
					Privilege:   "CreateManagedPrefixList",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a NAT gateway in a subnet",
					Privilege:   "CreateNatGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "elastic-ip*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "natgateway*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a network ACL in a VPC",
					Privilege:   "CreateNetworkAcl",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-acl*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a numbered entry (a rule) in a network ACL",
					Privilege:   "CreateNetworkAclEntry",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-acl*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a path to analyze for reachability",
					Privilege:   "CreateNetworkInsightsPath",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "network-insights-path*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a network interface in a subnet",
					Privilege:   "CreateNetworkInterface",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to create a permission for an AWS-authorized user to perform certain operations on a network interface",
					Privilege:   "CreateNetworkInterfacePermission",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a placement group",
					Privilege:   "CreatePlacementGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:PlacementGroupStrategy",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "placement-group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a root volume replacement task",
					Privilege:   "CreateReplaceRootVolumeTask",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "replace-root-volume-task*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a listing for Standard Reserved Instances to be sold in the Reserved Instance Marketplace",
					Privilege:   "CreateReservedInstancesListing",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:InstanceType",
								"ec2:Region",
								"ec2:ReservedInstancesOfferingType",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "reserved-instances*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to start a task that restores an AMI from an S3 object previously created by using CreateStoreImageTask",
					Privilege:   "CreateRestoreImageTask",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a route in a VPC route table",
					Privilege:   "CreateRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a route table for a VPC",
					Privilege:   "CreateRouteTable",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a security group",
					Privilege:   "CreateSecurityGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a snapshot of an EBS volume and store it in Amazon S3",
					Privilege:   "CreateSnapshot",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:OutpostArn",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:SnapshotTime",
								"ec2:SourceOutpostArn",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create crash-consistent snapshots of multiple EBS volumes and store them in Amazon S3",
					Privilege:   "CreateSnapshots",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:OutpostArn",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:SnapshotTime",
								"ec2:SourceOutpostArn",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a data feed for Spot Instances to view Spot Instance usage logs",
					Privilege:   "CreateSpotDatafeedSubscription",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to store an AMI as a single object in an S3 bucket",
					Privilege:   "CreateStoreImageTask",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a subnet in a VPC",
					Privilege:   "CreateSubnet",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a subnet CIDR reservation",
					Privilege:   "CreateSubnetCidrReservation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add or overwrite one or more tags for Amazon EC2 resources",
					Privilege:   "CreateTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "customer-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dhcp-options",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "egress-only-internet-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ElasticGpuType",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "elastic-gpu",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "elastic-ip",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "export-image-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "export-instance-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "fleet",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "fpga-image",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "host-reservation",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "import-image-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "import-snapshot-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "instance-event-window",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "internet-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "ipv4pool-ec2",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "ipv6pool-ec2",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table-virtual-interface-group-association",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table-vpc-association",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-virtual-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-virtual-interface-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "natgateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-acl",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:PlacementGroupStrategy",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "placement-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "replace-root-volume-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:InstanceType",
								"ec2:Region",
								"ec2:ReservedInstancesOfferingType",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "reserved-instances",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "security-group-rule",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "spot-fleet-request",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "spot-instances-request",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-session",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-target",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-connect-peer",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-flow-log",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AccepterVpc",
								"ec2:Region",
								"ec2:RequesterVpc",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-peering-connection",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway",
						},
						{
							ConditionKeys: []string{
								"ec2:CreateAction",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a traffic mirror filter",
					Privilege:   "CreateTrafficMirrorFilter",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a traffic mirror filter rule",
					Privilege:   "CreateTrafficMirrorFilterRule",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter*",
						},
						{
							ConditionKeys: []string{
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter-rule*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a traffic mirror session",
					Privilege:   "CreateTrafficMirrorSession",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-session*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-target*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a traffic mirror target",
					Privilege:   "CreateTrafficMirrorTarget",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-target*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a transit gateway",
					Privilege:   "CreateTransitGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a Connect attachment from a specified transit gateway attachment",
					Privilege:   "CreateTransitGatewayConnect",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a Connect peer between a transit gateway and an appliance",
					Privilege:   "CreateTransitGatewayConnectPeer",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a multicast domain for a transit gateway",
					Privilege:   "CreateTransitGatewayMulticastDomain",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to request a transit gateway peering attachment between a requester and accepter transit gateway",
					Privilege:   "CreateTransitGatewayPeeringAttachment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a transit gateway prefix list reference",
					Privilege:   "CreateTransitGatewayPrefixListReference",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a static route for a transit gateway route table",
					Privilege:   "CreateTransitGatewayRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a route table for a transit gateway",
					Privilege:   "CreateTransitGatewayRouteTable",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to attach a VPC to a transit gateway",
					Privilege:   "CreateTransitGatewayVpcAttachment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an EBS volume",
					Privilege:   "CreateVolume",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a VPC with a specified CIDR block",
					Privilege:   "CreateVpc",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "ipv6pool-ec2",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a VPC endpoint for an AWS service",
					Privilege:   "CreateVpcEndpoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{
								"route53:AssociateVPCWithHostedZone",
							},
							ResourceType: "vpc*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
								"ec2:VpceServiceName",
								"ec2:VpceServiceOwner",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a connection notification for a VPC endpoint or VPC endpoint service",
					Privilege:   "CreateVpcEndpointConnectionNotification",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a VPC endpoint service configuration to which service consumers (AWS accounts, IAM users, and IAM roles) can connect",
					Privilege:   "CreateVpcEndpointServiceConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to request a VPC peering connection between two VPCs",
					Privilege:   "CreateVpcPeeringConnection",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:AccepterVpc",
								"ec2:Region",
								"ec2:RequesterVpc",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-peering-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a VPN connection between a virtual private gateway or transit gateway and a customer gateway",
					Privilege:   "CreateVpnConnection",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "customer-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a static route for a VPN connection between a virtual private gateway and a customer gateway",
					Privilege:   "CreateVpnConnectionRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a virtual private gateway",
					Privilege:   "CreateVpnGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a carrier gateway",
					Privilege:   "DeleteCarrierGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "carrier-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a Client VPN endpoint",
					Privilege:   "DeleteClientVpnEndpoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a route from a Client VPN endpoint",
					Privilege:   "DeleteClientVpnRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a customer gateway",
					Privilege:   "DeleteCustomerGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "customer-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a set of DHCP options",
					Privilege:   "DeleteDhcpOptions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dhcp-options*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an egress-only internet gateway",
					Privilege:   "DeleteEgressOnlyInternetGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "egress-only-internet-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete one or more EC2 Fleets",
					Privilege:   "DeleteFleets",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "fleet*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete one or more flow logs",
					Privilege:   "DeleteFlowLogs",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-flow-log*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an Amazon FPGA Image (AFI)",
					Privilege:   "DeleteFpgaImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "fpga-image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified event window.",
					Privilege:   "DeleteInstanceEventWindow",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "instance-event-window*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an internet gateway",
					Privilege:   "DeleteInternetGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "internet-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a key pair by removing the public key from Amazon EC2",
					Privilege:   "DeleteKeyPair",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a launch template and its associated versions",
					Privilege:   "DeleteLaunchTemplate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete one or more versions of a launch template",
					Privilege:   "DeleteLaunchTemplateVersions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a route from a local gateway route table",
					Privilege:   "DeleteLocalGatewayRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an association between a VPC and local gateway route table",
					Privilege:   "DeleteLocalGatewayRouteTableVpcAssociation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table-vpc-association*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a managed prefix list",
					Privilege:   "DeleteManagedPrefixList",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a NAT gateway",
					Privilege:   "DeleteNatGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "natgateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a network ACL",
					Privilege:   "DeleteNetworkAcl",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-acl*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an inbound or outbound entry (rule) from a network ACL",
					Privilege:   "DeleteNetworkAclEntry",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-acl*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a network insights analysis",
					Privilege:   "DeleteNetworkInsightsAnalysis",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "network-insights-analysis*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a network insights path",
					Privilege:   "DeleteNetworkInsightsPath",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "network-insights-path*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a detached network interface",
					Privilege:   "DeleteNetworkInterface",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete a permission that is associated with a network interface",
					Privilege:   "DeleteNetworkInterfacePermission",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a placement group",
					Privilege:   "DeletePlacementGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the queued purchases for the specified Reserved Instances",
					Privilege:   "DeleteQueuedReservedInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:InstanceType",
								"ec2:Region",
								"ec2:ReservedInstancesOfferingType",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "reserved-instances*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a route from a route table",
					Privilege:   "DeleteRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a route table",
					Privilege:   "DeleteRouteTable",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a security group",
					Privilege:   "DeleteSecurityGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a snapshot of an EBS volume",
					Privilege:   "DeleteSnapshot",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:OutpostArn",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:SourceOutpostArn",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a data feed for Spot Instances",
					Privilege:   "DeleteSpotDatafeedSubscription",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a subnet",
					Privilege:   "DeleteSubnet",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a subnet CIDR reservation",
					Privilege:   "DeleteSubnetCidrReservation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to delete one or more tags from Amazon EC2 resources",
					Privilege:   "DeleteTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "customer-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dhcp-options",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "egress-only-internet-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ElasticGpuType",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "elastic-gpu",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "elastic-ip",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "export-image-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "export-instance-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "fleet",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "fpga-image",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "host-reservation",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "import-image-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "import-snapshot-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "instance-event-window",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "internet-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "ipv4pool-ec2",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "ipv6pool-ec2",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table-virtual-interface-group-association",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table-vpc-association",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-virtual-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-virtual-interface-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "natgateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-acl",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:PlacementGroupStrategy",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "placement-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "replace-root-volume-task",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:InstanceType",
								"ec2:Region",
								"ec2:ReservedInstancesOfferingType",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "reserved-instances",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "security-group-rule",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "spot-fleet-request",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "spot-instances-request",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-session",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-target",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-connect-peer",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-flow-log",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AccepterVpc",
								"ec2:Region",
								"ec2:RequesterVpc",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-peering-connection",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a traffic mirror filter",
					Privilege:   "DeleteTrafficMirrorFilter",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a traffic mirror filter rule",
					Privilege:   "DeleteTrafficMirrorFilterRule",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter*",
						},
						{
							ConditionKeys: []string{
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter-rule*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a traffic mirror session",
					Privilege:   "DeleteTrafficMirrorSession",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-session*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a traffic mirror target",
					Privilege:   "DeleteTrafficMirrorTarget",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-target*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a transit gateway",
					Privilege:   "DeleteTransitGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a transit gateway connect attachment",
					Privilege:   "DeleteTransitGatewayConnect",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a transit gateway connect peer",
					Privilege:   "DeleteTransitGatewayConnectPeer",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-connect-peer*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permissions to delete a transit gateway multicast domain",
					Privilege:   "DeleteTransitGatewayMulticastDomain",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a peering attachment from a transit gateway",
					Privilege:   "DeleteTransitGatewayPeeringAttachment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a transit gateway prefix list reference",
					Privilege:   "DeleteTransitGatewayPrefixListReference",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a route from a transit gateway route table",
					Privilege:   "DeleteTransitGatewayRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a transit gateway route table",
					Privilege:   "DeleteTransitGatewayRouteTable",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a VPC attachment from a transit gateway",
					Privilege:   "DeleteTransitGatewayVpcAttachment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an EBS volume",
					Privilege:   "DeleteVolume",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a VPC",
					Privilege:   "DeleteVpc",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete one or more VPC endpoint connection notifications",
					Privilege:   "DeleteVpcEndpointConnectionNotifications",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete one or more VPC endpoint service configurations",
					Privilege:   "DeleteVpcEndpointServiceConfigurations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete one or more VPC endpoints",
					Privilege:   "DeleteVpcEndpoints",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServiceName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a VPC peering connection",
					Privilege:   "DeleteVpcPeeringConnection",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AccepterVpc",
								"ec2:Region",
								"ec2:RequesterVpc",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-peering-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a VPN connection",
					Privilege:   "DeleteVpnConnection",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a static route for a VPN connection between a virtual private gateway and a customer gateway",
					Privilege:   "DeleteVpnConnectionRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a virtual private gateway",
					Privilege:   "DeleteVpnGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to release an IP address range that was provisioned through bring your own IP addresses (BYOIP), and to delete the corresponding address pool",
					Privilege:   "DeprovisionByoipCidr",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to deregister an Amazon Machine Image (AMI)",
					Privilege:   "DeregisterImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove tags from the set of tags to include in notifications about scheduled events for your instances",
					Privilege:   "DeregisterInstanceEventNotificationAttributes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to deregister one or more network interface members from a group IP address in a transit gateway multicast domain",
					Privilege:   "DeregisterTransitGatewayMulticastGroupMembers",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to deregister one or more network interface sources from a group IP address in a transit gateway multicast domain",
					Privilege:   "DeregisterTransitGatewayMulticastGroupSources",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the attributes of the AWS account",
					Privilege:   "DescribeAccountAttributes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more Elastic IP addresses",
					Privilege:   "DescribeAddresses",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the attributes of the specified Elastic IP addresses",
					Privilege:   "DescribeAddressesAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the longer ID format settings for all resource types",
					Privilege:   "DescribeAggregateIdFormat",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more of the Availability Zones that are available to you",
					Privilege:   "DescribeAvailabilityZones",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more bundling tasks",
					Privilege:   "DescribeBundleTasks",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the IP address ranges that were provisioned through bring your own IP addresses (BYOIP)",
					Privilege:   "DescribeByoipCidrs",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more Capacity Reservations",
					Privilege:   "DescribeCapacityReservations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more Carrier Gateways",
					Privilege:   "DescribeCarrierGateways",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more linked EC2-Classic instances",
					Privilege:   "DescribeClassicLinkInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the authorization rules for a Client VPN endpoint",
					Privilege:   "DescribeClientVpnAuthorizationRules",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe active client connections and connections that have been terminated within the last 60 minutes for a Client VPN endpoint",
					Privilege:   "DescribeClientVpnConnections",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more Client VPN endpoints",
					Privilege:   "DescribeClientVpnEndpoints",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the routes for a Client VPN endpoint",
					Privilege:   "DescribeClientVpnRoutes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the target networks that are associated with a Client VPN endpoint",
					Privilege:   "DescribeClientVpnTargetNetworks",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the specified customer-owned address pools or all of your customer-owned address pools",
					Privilege:   "DescribeCoipPools",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more conversion tasks",
					Privilege:   "DescribeConversionTasks",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more customer gateways",
					Privilege:   "DescribeCustomerGateways",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more DHCP options sets",
					Privilege:   "DescribeDhcpOptions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more egress-only internet gateways",
					Privilege:   "DescribeEgressOnlyInternetGateways",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe an Elastic Graphics accelerator that is associated with an instance",
					Privilege:   "DescribeElasticGpus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more export image tasks",
					Privilege:   "DescribeExportImageTasks",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more export instance tasks",
					Privilege:   "DescribeExportTasks",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe the state of fast snapshot restores for snapshots",
					Privilege:   "DescribeFastSnapshotRestores",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the events for an EC2 Fleet during a specified time",
					Privilege:   "DescribeFleetHistory",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the running instances for an EC2 Fleet",
					Privilege:   "DescribeFleetInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more EC2 Fleets",
					Privilege:   "DescribeFleets",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more flow logs",
					Privilege:   "DescribeFlowLogs",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the attributes of an Amazon FPGA Image (AFI)",
					Privilege:   "DescribeFpgaImageAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more Amazon FPGA Images (AFIs)",
					Privilege:   "DescribeFpgaImages",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the Dedicated Host Reservations that are available to purchase",
					Privilege:   "DescribeHostReservationOfferings",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the Dedicated Host Reservations that are associated with Dedicated Hosts in the AWS account",
					Privilege:   "DescribeHostReservations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more Dedicated Hosts",
					Privilege:   "DescribeHosts",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the IAM instance profile associations",
					Privilege:   "DescribeIamInstanceProfileAssociations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the ID format settings for resources",
					Privilege:   "DescribeIdFormat",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the ID format settings for resources for an IAM user, IAM role, or root user",
					Privilege:   "DescribeIdentityIdFormat",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe an attribute of an Amazon Machine Image (AMI)",
					Privilege:   "DescribeImageAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more images (AMIs, AKIs, and ARIs)",
					Privilege:   "DescribeImages",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe import virtual machine or import snapshot tasks",
					Privilege:   "DescribeImportImageTasks",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe import snapshot tasks",
					Privilege:   "DescribeImportSnapshotTasks",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the attributes of an instance",
					Privilege:   "DescribeInstanceAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the credit option for CPU usage of one or more burstable performance instances",
					Privilege:   "DescribeInstanceCreditSpecifications",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the set of tags to include in notifications about scheduled events for your instances",
					Privilege:   "DescribeInstanceEventNotificationAttributes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the specified event windows or all event windows",
					Privilege:   "DescribeInstanceEventWindows",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the status of one or more instances",
					Privilege:   "DescribeInstanceStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the set of instance types that are offered in a location",
					Privilege:   "DescribeInstanceTypeOfferings",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the details of instance types that are offered in a location",
					Privilege:   "DescribeInstanceTypes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more instances",
					Privilege:   "DescribeInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more internet gateways",
					Privilege:   "DescribeInternetGateways",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more IPv6 address pools",
					Privilege:   "DescribeIpv6Pools",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more key pairs",
					Privilege:   "DescribeKeyPairs",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more launch template versions",
					Privilege:   "DescribeLaunchTemplateVersions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more launch templates",
					Privilege:   "DescribeLaunchTemplates",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the associations between virtual interface groups and local gateway route tables",
					Privilege:   "DescribeLocalGatewayRouteTableVirtualInterfaceGroupAssociations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe an association between VPCs and local gateway route tables",
					Privilege:   "DescribeLocalGatewayRouteTableVpcAssociations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more local gateway route tables",
					Privilege:   "DescribeLocalGatewayRouteTables",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe local gateway virtual interface groups",
					Privilege:   "DescribeLocalGatewayVirtualInterfaceGroups",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe local gateway virtual interfaces",
					Privilege:   "DescribeLocalGatewayVirtualInterfaces",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more local gateways",
					Privilege:   "DescribeLocalGateways",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe your managed prefix lists and any AWS-managed prefix lists",
					Privilege:   "DescribeManagedPrefixLists",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe Elastic IP addresses that are being moved to the EC2-VPC platform",
					Privilege:   "DescribeMovingAddresses",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more NAT gateways",
					Privilege:   "DescribeNatGateways",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more network ACLs",
					Privilege:   "DescribeNetworkAcls",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more network insights analyses",
					Privilege:   "DescribeNetworkInsightsAnalyses",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more network insights paths",
					Privilege:   "DescribeNetworkInsightsPaths",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe a network interface attribute",
					Privilege:   "DescribeNetworkInterfaceAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the permissions that are associated with a network interface",
					Privilege:   "DescribeNetworkInterfacePermissions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more network interfaces",
					Privilege:   "DescribeNetworkInterfaces",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more placement groups",
					Privilege:   "DescribePlacementGroups",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe available AWS services in a prefix list format",
					Privilege:   "DescribePrefixLists",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the ID format settings for the root user and all IAM roles and IAM users that have explicitly specified a longer ID (17-character ID) preference",
					Privilege:   "DescribePrincipalIdFormat",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more IPv4 address pools",
					Privilege:   "DescribePublicIpv4Pools",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more AWS Regions that are currently available in your account",
					Privilege:   "DescribeRegions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe a root volume replacement task",
					Privilege:   "DescribeReplaceRootVolumeTasks",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more purchased Reserved Instances in your account",
					Privilege:   "DescribeReservedInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe your account's Reserved Instance listings in the Reserved Instance Marketplace",
					Privilege:   "DescribeReservedInstancesListings",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the modifications made to one or more Reserved Instances",
					Privilege:   "DescribeReservedInstancesModifications",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the Reserved Instance offerings that are available for purchase",
					Privilege:   "DescribeReservedInstancesOfferings",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more route tables",
					Privilege:   "DescribeRouteTables",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to find available schedules for Scheduled Instances",
					Privilege:   "DescribeScheduledInstanceAvailability",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe one or more Scheduled Instances in your account",
					Privilege:   "DescribeScheduledInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the VPCs on the other side of a VPC peering connection that are referencing specified VPC security groups",
					Privilege:   "DescribeSecurityGroupReferences",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more of your security group rules",
					Privilege:   "DescribeSecurityGroupRules",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more security groups",
					Privilege:   "DescribeSecurityGroups",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe an attribute of a snapshot",
					Privilege:   "DescribeSnapshotAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more EBS snapshots",
					Privilege:   "DescribeSnapshots",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the data feed for Spot Instances",
					Privilege:   "DescribeSpotDatafeedSubscription",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the running instances for a Spot Fleet",
					Privilege:   "DescribeSpotFleetInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the events for a Spot Fleet request during a specified time",
					Privilege:   "DescribeSpotFleetRequestHistory",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more Spot Fleet requests",
					Privilege:   "DescribeSpotFleetRequests",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more Spot Instance requests",
					Privilege:   "DescribeSpotInstanceRequests",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the Spot Instance price history",
					Privilege:   "DescribeSpotPriceHistory",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the stale security group rules for security groups in a specified VPC",
					Privilege:   "DescribeStaleSecurityGroups",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the progress of the AMI store tasks",
					Privilege:   "DescribeStoreImageTasks",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more subnets",
					Privilege:   "DescribeSubnets",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe one or more tags for an Amazon EC2 resource",
					Privilege:   "DescribeTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more traffic mirror filters",
					Privilege:   "DescribeTrafficMirrorFilters",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more traffic mirror sessions",
					Privilege:   "DescribeTrafficMirrorSessions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more traffic mirror targets",
					Privilege:   "DescribeTrafficMirrorTargets",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more attachments between resources and transit gateways",
					Privilege:   "DescribeTransitGatewayAttachments",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more transit gateway connect peers",
					Privilege:   "DescribeTransitGatewayConnectPeers",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more transit gateway connect attachments",
					Privilege:   "DescribeTransitGatewayConnects",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more transit gateway multicast domains",
					Privilege:   "DescribeTransitGatewayMulticastDomains",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more transit gateway peering attachments",
					Privilege:   "DescribeTransitGatewayPeeringAttachments",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more transit gateway route tables",
					Privilege:   "DescribeTransitGatewayRouteTables",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more VPC attachments on a transit gateway",
					Privilege:   "DescribeTransitGatewayVpcAttachments",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more transit gateways",
					Privilege:   "DescribeTransitGateways",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more network interface trunk associations",
					Privilege:   "DescribeTrunkInterfaceAssociations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe an attribute of an EBS volume",
					Privilege:   "DescribeVolumeAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the status of one or more EBS volumes",
					Privilege:   "DescribeVolumeStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more EBS volumes",
					Privilege:   "DescribeVolumes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe the current modification status of one or more EBS volumes",
					Privilege:   "DescribeVolumesModifications",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe an attribute of a VPC",
					Privilege:   "DescribeVpcAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the ClassicLink status of one or more VPCs",
					Privilege:   "DescribeVpcClassicLink",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the ClassicLink DNS support status of one or more VPCs",
					Privilege:   "DescribeVpcClassicLinkDnsSupport",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the connection notifications for VPC endpoints and VPC endpoint services",
					Privilege:   "DescribeVpcEndpointConnectionNotifications",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the VPC endpoint connections to your VPC endpoint services",
					Privilege:   "DescribeVpcEndpointConnections",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe VPC endpoint service configurations (your services)",
					Privilege:   "DescribeVpcEndpointServiceConfigurations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe the principals (service consumers) that are permitted to discover your VPC endpoint service",
					Privilege:   "DescribeVpcEndpointServicePermissions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe all supported AWS services that can be specified when creating a VPC endpoint",
					Privilege:   "DescribeVpcEndpointServices",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more VPC endpoints",
					Privilege:   "DescribeVpcEndpoints",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more VPC peering connections",
					Privilege:   "DescribeVpcPeeringConnections",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more VPCs",
					Privilege:   "DescribeVpcs",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe one or more VPN connections",
					Privilege:   "DescribeVpnConnections",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe one or more virtual private gateways",
					Privilege:   "DescribeVpnGateways",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to unlink (detach) a linked EC2-Classic instance from a VPC",
					Privilege:   "DetachClassicLinkVpc",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to detach an internet gateway from a VPC",
					Privilege:   "DetachInternetGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "internet-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to detach a network interface from an instance",
					Privilege:   "DetachNetworkInterface",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to detach an EBS volume from an instance",
					Privilege:   "DetachVolume",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to detach a virtual private gateway from a VPC",
					Privilege:   "DetachVpnGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disable EBS encryption by default for your account",
					Privilege:   "DisableEbsEncryptionByDefault",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disable fast snapshot restores for one or more snapshots in specified Availability Zones",
					Privilege:   "DisableFastSnapshotRestores",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to cancel the deprecation of the specified AMI",
					Privilege:   "DisableImageDeprecation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disable access to the EC2 serial console of all instances for your account",
					Privilege:   "DisableSerialConsoleAccess",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disable a resource attachment from propagating routes to the specified propagation route table",
					Privilege:   "DisableTransitGatewayRouteTablePropagation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disable a virtual private gateway from propagating routes to a specified route table of a VPC",
					Privilege:   "DisableVgwRoutePropagation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disable ClassicLink for a VPC",
					Privilege:   "DisableVpcClassicLink",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disable ClassicLink DNS support for a VPC",
					Privilege:   "DisableVpcClassicLinkDnsSupport",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate an Elastic IP address from an instance or network interface",
					Privilege:   "DisassociateAddress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "elastic-ip",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate a target network from a Client VPN endpoint",
					Privilege:   "DisassociateClientVpnTargetNetwork",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate an ACM certificate from a IAM role",
					Privilege:   "DisassociateEnclaveCertificateIamRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "certificate*",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate an IAM instance profile from a running or stopped instance",
					Privilege:   "DisassociateIamInstanceProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate one or more targets from an event window",
					Privilege:   "DisassociateInstanceEventWindow",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "instance-event-window*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate a subnet from a route table",
					Privilege:   "DisassociateRouteTable",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate a CIDR block from a subnet",
					Privilege:   "DisassociateSubnetCidrBlock",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate one or more subnets from a transit gateway multicast domain",
					Privilege:   "DisassociateTransitGatewayMulticastDomain",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate a resource attachment from a transit gateway route table",
					Privilege:   "DisassociateTransitGatewayRouteTable",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate a branch network interface to a trunk network interface",
					Privilege:   "DisassociateTrunkInterface",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disassociate a CIDR block from a VPC",
					Privilege:   "DisassociateVpcCidrBlock",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable EBS encryption by default for your account",
					Privilege:   "EnableEbsEncryptionByDefault",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable fast snapshot restores for one or more snapshots in specified Availability Zones",
					Privilege:   "EnableFastSnapshotRestores",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable deprecation of the specified AMI at the specified date and time",
					Privilege:   "EnableImageDeprecation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable access to the EC2 serial console of all instances for your account",
					Privilege:   "EnableSerialConsoleAccess",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable an attachment to propagate routes to a propagation route table",
					Privilege:   "EnableTransitGatewayRouteTablePropagation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable a virtual private gateway to propagate routes to a VPC route table",
					Privilege:   "EnableVgwRoutePropagation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable I/O operations for a volume that had I/O operations disabled",
					Privilege:   "EnableVolumeIO",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable a VPC for ClassicLink",
					Privilege:   "EnableVpcClassicLink",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable a VPC to support DNS hostname resolution for ClassicLink",
					Privilege:   "EnableVpcClassicLinkDnsSupport",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to download the client certificate revocation list for a Client VPN endpoint",
					Privilege:   "ExportClientVpnClientCertificateRevocationList",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to download the contents of the Client VPN endpoint configuration file for a Client VPN endpoint",
					Privilege:   "ExportClientVpnClientConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to export an Amazon Machine Image (AMI) to a VM file",
					Privilege:   "ExportImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "export-image-task*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to export routes from a transit gateway route table to an Amazon S3 bucket",
					Privilege:   "ExportTransitGatewayRoutes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get the list of roles associated with an ACM certificate",
					Privilege:   "GetAssociatedEnclaveCertificateIamRoles",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "certificate*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get information about the IPv6 CIDR block associations for a specified IPv6 address pool",
					Privilege:   "GetAssociatedIpv6PoolCidrs",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "ipv6pool-ec2*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get usage information about a Capacity Reservation",
					Privilege:   "GetCapacityReservationUsage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe the allocations from the specified customer-owned address pool",
					Privilege:   "GetCoipPoolUsage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get the console output for an instance",
					Privilege:   "GetConsoleOutput",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve a JPG-format screenshot of a running instance",
					Privilege:   "GetConsoleScreenshot",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get the default credit option for CPU usage of a burstable performance instance family",
					Privilege:   "GetDefaultCreditSpecification",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get the ID of the default customer master key (CMK) for EBS encryption by default",
					Privilege:   "GetEbsDefaultKmsKeyId",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe whether EBS encryption by default is enabled for your account",
					Privilege:   "GetEbsEncryptionByDefault",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to generate a CloudFormation template to streamline the integration of VPC flow logs with Amazon Athena",
					Privilege:   "GetFlowLogsIntegrationTemplate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-flow-log*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the resource groups to which a Capacity Reservation has been added",
					Privilege:   "GetGroupsForCapacityReservation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to preview a reservation purchase with configurations that match those of a Dedicated Host",
					Privilege:   "GetHostReservationPurchasePreview",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get the configuration data of the specified instance for use with a new launch template or launch template version",
					Privilege:   "GetLaunchTemplateData",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get information about the resources that are associated with the specified managed prefix list",
					Privilege:   "GetManagedPrefixListAssociations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get information about the entries for a specified managed prefix list",
					Privilege:   "GetManagedPrefixListEntries",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the encrypted administrator password for a running Windows instance",
					Privilege:   "GetPasswordData",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return a quote and exchange information for exchanging one or more Convertible Reserved Instances for a new Convertible Reserved Instance",
					Privilege:   "GetReservedInstancesExchangeQuote",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:InstanceType",
								"ec2:Region",
								"ec2:ReservedInstancesOfferingType",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "reserved-instances",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the access status of your account to the EC2 serial console of all instances",
					Privilege:   "GetSerialConsoleAccessStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about the subnet CIDR reservations",
					Privilege:   "GetSubnetCidrReservations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the route tables to which a resource attachment propagates routes",
					Privilege:   "GetTransitGatewayAttachmentPropagations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to get information about the associations for a transit gateway multicast domain",
					Privilege:   "GetTransitGatewayMulticastDomainAssociations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to get information about prefix list references for a transit gateway route table",
					Privilege:   "GetTransitGatewayPrefixListReferences",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to get information about associations for a transit gateway route table",
					Privilege:   "GetTransitGatewayRouteTableAssociations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to get information about the route table propagations for a transit gateway route table",
					Privilege:   "GetTransitGatewayRouteTablePropagations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to upload a client certificate revocation list to a Client VPN endpoint",
					Privilege:   "ImportClientVpnClientCertificateRevocationList",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to import single or multi-volume disk images or EBS snapshots into an Amazon Machine Image (AMI)",
					Privilege:   "ImportImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an import instance task using metadata from a disk image",
					Privilege:   "ImportInstance",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to import a public key from an RSA key pair that was created with a third-party tool",
					Privilege:   "ImportKeyPair",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:KeyPairName",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to import a disk into an EBS snapshot",
					Privilege:   "ImportSnapshot",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an import volume task using metadata from a disk image",
					Privilege:   "ImportVolume",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an attribute of the specified Elastic IP address",
					Privilege:   "ModifyAddressAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the opt-in status of the Local Zone and Wavelength Zone group for your account",
					Privilege:   "ModifyAvailabilityZoneGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a Capacity Reservation's capacity and the conditions under which it is to be released",
					Privilege:   "ModifyCapacityReservation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:CapacityReservationFleet",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a Client VPN endpoint",
					Privilege:   "ModifyClientVpnEndpoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to change the account level default credit option for CPU usage of burstable performance instances",
					Privilege:   "ModifyDefaultCreditSpecification",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to change the default customer master key (CMK) for EBS encryption by default for your account",
					Privilege:   "ModifyEbsDefaultKmsKeyId",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an EC2 Fleet",
					Privilege:   "ModifyFleet",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "fleet*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an attribute of an Amazon FPGA Image (AFI)",
					Privilege:   "ModifyFpgaImageAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "fpga-image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a Dedicated Host",
					Privilege:   "ModifyHosts",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the ID format for a resource",
					Privilege:   "ModifyIdFormat",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the ID format of a resource for a specific principal in your account",
					Privilege:   "ModifyIdentityIdFormat",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an attribute of an Amazon Machine Image (AMI)",
					Privilege:   "ModifyImageAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an attribute of an instance",
					Privilege:   "ModifyInstanceAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the Capacity Reservation settings for a stopped instance",
					Privilege:   "ModifyInstanceCapacityReservationAttributes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the credit option for CPU usage on an instance",
					Privilege:   "ModifyInstanceCreditSpecification",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the start time for a scheduled EC2 instance event",
					Privilege:   "ModifyInstanceEventStartTime",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the specified event window",
					Privilege:   "ModifyInstanceEventWindow",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "instance-event-window*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the metadata options for an instance",
					Privilege:   "ModifyInstanceMetadataOptions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the placement attributes for an instance",
					Privilege:   "ModifyInstancePlacement",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:PlacementGroupStrategy",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "placement-group",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a launch template",
					Privilege:   "ModifyLaunchTemplate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a managed prefix list",
					Privilege:   "ModifyManagedPrefixList",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an attribute of a network interface",
					Privilege:   "ModifyNetworkInterfaceAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:Attribute/${AttributeName}",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify attributes of one or more Reserved Instances",
					Privilege:   "ModifyReservedInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:InstanceType",
								"ec2:Region",
								"ec2:ReservedInstancesOfferingType",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "reserved-instances*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the rules of a security group",
					Privilege:   "ModifySecurityGroupRules",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "security-group-rule",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to add or remove permission settings for a snapshot",
					Privilege:   "ModifySnapshotAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a Spot Fleet request",
					Privilege:   "ModifySpotFleetRequest",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "spot-fleet-request*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an attribute of a subnet",
					Privilege:   "ModifySubnetAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to allow or restrict mirroring network services",
					Privilege:   "ModifyTrafficMirrorFilterNetworkServices",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a traffic mirror rule",
					Privilege:   "ModifyTrafficMirrorFilterRule",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter*",
						},
						{
							ConditionKeys: []string{
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter-rule*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a traffic mirror session",
					Privilege:   "ModifyTrafficMirrorSession",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-session*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-filter",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "traffic-mirror-target",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a transit gateway",
					Privilege:   "ModifyTransitGateway",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a transit gateway prefix list reference",
					Privilege:   "ModifyTransitGatewayPrefixListReference",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a VPC attachment on a transit gateway",
					Privilege:   "ModifyTransitGatewayVpcAttachment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the parameters of an EBS volume",
					Privilege:   "ModifyVolume",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an attribute of a volume",
					Privilege:   "ModifyVolumeAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an attribute of a VPC",
					Privilege:   "ModifyVpcAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify an attribute of a VPC endpoint",
					Privilege:   "ModifyVpcEndpoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify a connection notification for a VPC endpoint or VPC endpoint service",
					Privilege:   "ModifyVpcEndpointConnectionNotification",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the attributes of a VPC endpoint service configuration",
					Privilege:   "ModifyVpcEndpointServiceConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to modify the permissions for a VPC endpoint service",
					Privilege:   "ModifyVpcEndpointServicePermissions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the VPC peering connection options on one side of a VPC peering connection",
					Privilege:   "ModifyVpcPeeringConnectionOptions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AccepterVpc",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:RequesterVpc",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-peering-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the instance tenancy attribute of a VPC",
					Privilege:   "ModifyVpcTenancy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "vpc*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the target gateway of a Site-to-Site VPN connection",
					Privilege:   "ModifyVpnConnection",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "customer-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the connection options for your Site-to-Site VPN connection",
					Privilege:   "ModifyVpnConnectionOptions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the certificate for a Site-to-Site VPN connection",
					Privilege:   "ModifyVpnTunnelCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to modify the options for a Site-to-Site VPN connection",
					Privilege:   "ModifyVpnTunnelOptions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Attribute/${AttributeName}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable detailed monitoring for a running instance",
					Privilege:   "MonitorInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to move an Elastic IP address from the EC2-Classic platform to the EC2-VPC platform",
					Privilege:   "MoveAddressToVpc",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to provision an address range for use in AWS through bring your own IP addresses (BYOIP), and to create a corresponding address pool",
					Privilege:   "ProvisionByoipCidr",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to purchase a reservation with configurations that match those of a Dedicated Host",
					Privilege:   "PurchaseHostReservation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to purchase a Reserved Instance offering",
					Privilege:   "PurchaseReservedInstancesOffering",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:InstanceType",
								"ec2:Region",
								"ec2:ReservedInstancesOfferingType",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "reserved-instances*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to purchase one or more Scheduled Instances with a specified schedule",
					Privilege:   "PurchaseScheduledInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to request a reboot of one or more instances",
					Privilege:   "RebootInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to register an Amazon Machine Image (AMI)",
					Privilege:   "RegisterImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add tags to the set of tags to include in notifications about scheduled events for your instances",
					Privilege:   "RegisterInstanceEventNotificationAttributes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to register one or more network interfaces as a member of a group IP address in a transit gateway multicast domain",
					Privilege:   "RegisterTransitGatewayMulticastGroupMembers",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to register one or more network interfaces as a source of a group IP address in a transit gateway multicast domain",
					Privilege:   "RegisterTransitGatewayMulticastGroupSources",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reject requests to associate cross-account subnets with a transit gateway multicast domain",
					Privilege:   "RejectTransitGatewayMulticastDomainAssociations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reject a transit gateway peering attachment request",
					Privilege:   "RejectTransitGatewayPeeringAttachment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reject a request to attach a VPC to a transit gateway",
					Privilege:   "RejectTransitGatewayVpcAttachment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reject one or more VPC endpoint connection requests to a VPC endpoint service",
					Privilege:   "RejectVpcEndpointConnections",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reject a VPC peering connection request",
					Privilege:   "RejectVpcPeeringConnection",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AccepterVpc",
								"ec2:Region",
								"ec2:RequesterVpc",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-peering-connection*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to release an Elastic IP address",
					Privilege:   "ReleaseAddress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "elastic-ip",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to release one or more On-Demand Dedicated Hosts",
					Privilege:   "ReleaseHosts",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AutoPlacement",
								"ec2:AvailabilityZone",
								"ec2:HostRecovery",
								"ec2:InstanceType",
								"ec2:Quantity",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "dedicated-host*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to replace an IAM instance profile for an instance",
					Privilege:   "ReplaceIamInstanceProfileAssociation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:NewInstanceProfile",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{
								"iam:PassRole",
							},
							ResourceType: "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to change which network ACL a subnet is associated with",
					Privilege:   "ReplaceNetworkAclAssociation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-acl*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to replace an entry (rule) in a network ACL",
					Privilege:   "ReplaceNetworkAclEntry",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-acl*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to replace a route within a route table in a VPC",
					Privilege:   "ReplaceRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Tenancy",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "carrier-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "egress-only-internet-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "internet-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "natgateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AccepterVpc",
								"ec2:Region",
								"ec2:RequesterVpc",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-peering-connection",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-gateway",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to change the route table that is associated with a subnet",
					Privilege:   "ReplaceRouteTableAssociation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to replace a route in a transit gateway route table",
					Privilege:   "ReplaceTransitGatewayRoute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-attachment",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to submit feedback about the status of an instance",
					Privilege:   "ReportInstanceStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a Spot Fleet request",
					Privilege:   "RequestSpotFleet",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a Spot Instance request",
					Privilege:   "RequestSpotInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:Region",
							},
							DependentActions: []string{},
							ResourceType:     "spot-instances-request*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reset the attribute of the specified IP address",
					Privilege:   "ResetAddressAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reset the default customer master key (CMK) for EBS encryption for your account to use the AWS-managed CMK for EBS",
					Privilege:   "ResetEbsDefaultKmsKeyId",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reset an attribute of an Amazon FPGA Image (AFI) to its default value",
					Privilege:   "ResetFpgaImageAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "fpga-image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reset an attribute of an Amazon Machine Image (AMI) to its default value",
					Privilege:   "ResetImageAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reset an attribute of an instance to its default value",
					Privilege:   "ResetInstanceAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reset an attribute of a network interface",
					Privilege:   "ResetNetworkInterfaceAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to reset permission settings for a snapshot",
					Privilege:   "ResetSnapshotAttribute",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to restore an Elastic IP address that was previously moved to the EC2-VPC platform back to the EC2-Classic platform",
					Privilege:   "RestoreAddressToClassic",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to restore the entries from a previous version of a managed prefix list to a new version of the prefix list",
					Privilege:   "RestoreManagedPrefixListVersion",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "prefix-list*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove an inbound authorization rule from a Client VPN endpoint",
					Privilege:   "RevokeClientVpnIngress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove one or more outbound rules from a VPC security group",
					Privilege:   "RevokeSecurityGroupEgress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove one or more inbound rules from a security group",
					Privilege:   "RevokeSecurityGroupIngress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to launch one or more instances",
					Privilege:   "RunInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
								"ec2:AvailabilityZone",
								"ec2:Encrypted",
								"ec2:ParentSnapshot",
								"ec2:Region",
								"ec2:VolumeIops",
								"ec2:VolumeSize",
								"ec2:VolumeThroughput",
								"ec2:VolumeType",
							},
							DependentActions: []string{},
							ResourceType:     "volume*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "capacity-reservation",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ElasticGpuType",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "elastic-gpu",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "elastic-inference",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "launch-template",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:PlacementGroupStrategy",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "placement-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to launch one or more Scheduled Instances",
					Privilege:   "RunScheduledInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ImageType",
								"ec2:Owner",
								"ec2:Public",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
							},
							DependentActions: []string{},
							ResourceType:     "image*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:KeyPairName",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "key-pair",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:PlacementGroupStrategy",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "placement-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Owner",
								"ec2:ParentVolume",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SnapshotTime",
								"ec2:VolumeSize",
							},
							DependentActions: []string{},
							ResourceType:     "snapshot",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "subnet",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to search for routes in a local gateway route table",
					Privilege:   "SearchLocalGatewayRoutes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "local-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to search for groups, sources, and members in a transit gateway multicast domain",
					Privilege:   "SearchTransitGatewayMulticastGroups",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-multicast-domain",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to search for routes in a transit gateway route table",
					Privilege:   "SearchTransitGatewayRoutes",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "transit-gateway-route-table*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to send a diagnostic interrupt to an Amazon EC2 instance",
					Privilege:   "SendDiagnosticInterrupt",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to start a stopped instance",
					Privilege:   "StartInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to start analyzing a specified path",
					Privilege:   "StartNetworkInsightsAnalysis",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "network-insights-analysis*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "network-insights-path*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to start the private DNS verification process for a VPC endpoint service",
					Privilege:   "StartVpcEndpointServicePrivateDnsVerification",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:VpceServicePrivateDnsName",
							},
							DependentActions: []string{},
							ResourceType:     "vpc-endpoint-service*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to stop an Amazon EBS-backed instance",
					Privilege:   "StopInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to terminate active Client VPN endpoint connections",
					Privilege:   "TerminateClientVpnConnections",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:ClientRootCertificateChainArn",
								"ec2:CloudwatchLogGroupArn",
								"ec2:CloudwatchLogStreamArn",
								"ec2:DirectoryArn",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:SamlProviderArn",
								"ec2:ServerCertificateArn",
							},
							DependentActions: []string{},
							ResourceType:     "client-vpn-endpoint*",
						},
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AuthenticationType",
								"ec2:DPDTimeoutSeconds",
								"ec2:GatewayType",
								"ec2:IKEVersions",
								"ec2:InsideTunnelCidr",
								"ec2:Phase1DHGroupNumbers",
								"ec2:Phase1EncryptionAlgorithms",
								"ec2:Phase1IntegrityAlgorithms",
								"ec2:Phase1LifetimeSeconds",
								"ec2:Phase2DHGroupNumbers",
								"ec2:Phase2EncryptionAlgorithms",
								"ec2:Phase2IntegrityAlgorithms",
								"ec2:Phase2LifetimeSeconds",
								"ec2:PresharedKeys",
								"ec2:Region",
								"ec2:RekeyFuzzPercentage",
								"ec2:RekeyMarginTimeSeconds",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RoutingType",
							},
							DependentActions: []string{},
							ResourceType:     "vpn-connection",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to shut down one or more instances",
					Privilege:   "TerminateInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to unassign one or more IPv6 addresses from a network interface",
					Privilege:   "UnassignIpv6Addresses",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to unassign one or more secondary private IP addresses from a network interface",
					Privilege:   "UnassignPrivateIpAddresses",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AssociatePublicIpAddress",
								"ec2:AuthorizedService",
								"ec2:AvailabilityZone",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Subnet",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "network-interface*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to disable detailed monitoring for a running instance",
					Privilege:   "UnmonitorInstances",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:AvailabilityZone",
								"ec2:EbsOptimized",
								"ec2:InstanceProfile",
								"ec2:InstanceType",
								"ec2:PlacementGroup",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:RootDeviceType",
								"ec2:Tenancy",
							},
							DependentActions: []string{},
							ResourceType:     "instance*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update descriptions for one or more outbound rules in a VPC security group",
					Privilege:   "UpdateSecurityGroupRuleDescriptionsEgress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update descriptions for one or more inbound rules in a security group",
					Privilege:   "UpdateSecurityGroupRuleDescriptionsIngress",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:ResourceTag/${TagKey}",
								"ec2:Region",
								"ec2:ResourceTag/${TagKey}",
								"ec2:Vpc",
							},
							DependentActions: []string{},
							ResourceType:     "security-group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to stop advertising an address range that was provisioned for use in AWS through bring your own IP addresses (BYOIP)",
					Privilege:   "WithdrawByoipCidr",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
			},
			Resources: []ParliamentResource{
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:elastic-ip/${AllocationId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "elastic-ip",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:capacity-reservation/${CapacityReservationId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "capacity-reservation",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:carrier-gateway/${CarrierGatewayId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:Tenancy",
						"ec2:Vpc",
					},
					Resource: "carrier-gateway",
				},
				{
					Arn:           "arn:${Partition}:acm:${Region}:${Account}:certificate/${CertificateId}",
					ConditionKeys: []string{},
					Resource:      "certificate",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:client-vpn-endpoint/${ClientVpnEndpointId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:ClientRootCertificateChainArn",
						"ec2:CloudwatchLogGroupArn",
						"ec2:CloudwatchLogStreamArn",
						"ec2:DirectoryArn",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:SamlProviderArn",
						"ec2:ServerCertificateArn",
					},
					Resource: "client-vpn-endpoint",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:customer-gateway/${CustomerGatewayId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "customer-gateway",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:dedicated-host/${DedicatedHostId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:AutoPlacement",
						"ec2:AvailabilityZone",
						"ec2:HostRecovery",
						"ec2:InstanceType",
						"ec2:Quantity",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "dedicated-host",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:dhcp-options/${DhcpOptionsId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "dhcp-options",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:egress-only-internet-gateway/${EgressOnlyInternetGatewayId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "egress-only-internet-gateway",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:elastic-gpu/${ElasticGpuId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:ElasticGpuType",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "elastic-gpu",
				},
				{
					Arn:           "arn:${Partition}:elastic-inference:${Region}:${Account}:elastic-inference-accelerator/${ElasticInferenceAcceleratorId}",
					ConditionKeys: []string{},
					Resource:      "elastic-inference",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:export-image-task/${ExportImageTaskId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "export-image-task",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:export-instance-task/${ExportTaskId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "export-instance-task",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:fleet/${FleetId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "fleet",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}::fpga-image/${FpgaImageId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Owner",
						"ec2:Public",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "fpga-image",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:host-reservation/${HostReservationId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "host-reservation",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}::image/${ImageId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:ImageType",
						"ec2:Owner",
						"ec2:Public",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:RootDeviceType",
					},
					Resource: "image",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:import-image-task/${ImportImageTaskId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "import-image-task",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:import-snapshot-task/${ImportSnapshotTaskId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "import-snapshot-task",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:instance-event-window/${InstanceEventWindowId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "instance-event-window",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:instance/${InstanceId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:AvailabilityZone",
						"ec2:EbsOptimized",
						"ec2:InstanceProfile",
						"ec2:InstanceType",
						"ec2:PlacementGroup",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:RootDeviceType",
						"ec2:Tenancy",
					},
					Resource: "instance",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:internet-gateway/${InternetGatewayId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "internet-gateway",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:ipv4pool-ec2/${Ipv4PoolEc2Id}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "ipv4pool-ec2",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:ipv6pool-ec2/${Ipv6PoolEc2Id}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "ipv6pool-ec2",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:key-pair/${KeyPairName}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:KeyPairName",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "key-pair",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:launch-template/${LaunchTemplateId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "launch-template",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:local-gateway/${LocalGatewayId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "local-gateway",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:local-gateway-route-table-virtual-interface-group-association/${LocalGatewayRouteTableVirtualInterfaceGroupAssociationId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "local-gateway-route-table-virtual-interface-group-association",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:local-gateway-route-table-vpc-association/${LocalGatewayRouteTableVpcAssociationId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:Tenancy",
					},
					Resource: "local-gateway-route-table-vpc-association",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:local-gateway-route-table/${LocalGatewayRoutetableId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "local-gateway-route-table",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:local-gateway-virtual-interface-group/${LocalGatewayVirtualInterfaceGroupId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "local-gateway-virtual-interface-group",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:local-gateway-virtual-interface/${LocalGatewayVirtualInterfaceId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "local-gateway-virtual-interface",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:natgateway/${NatGatewayId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "natgateway",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:network-acl/${NaclId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:Vpc",
					},
					Resource: "network-acl",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:network-insights-analysis/${NetworkInsightsAnalysisId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "network-insights-analysis",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:network-insights-path/${NetworkInsightsPathId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "network-insights-path",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:network-interface/${NetworkInterfaceId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:AssociatePublicIpAddress",
						"ec2:AuthorizedService",
						"ec2:AvailabilityZone",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:Subnet",
						"ec2:Vpc",
					},
					Resource: "network-interface",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:placement-group/${PlacementGroupName}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:PlacementGroupStrategy",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "placement-group",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:prefix-list/${PrefixListId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "prefix-list",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:replace-root-volume-task/${ReplaceRootVolumeTaskId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "replace-root-volume-task",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:reserved-instances/${ReservationId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:AvailabilityZone",
						"ec2:InstanceType",
						"ec2:Region",
						"ec2:ReservedInstancesOfferingType",
						"ec2:ResourceTag/${TagKey}",
						"ec2:Tenancy",
					},
					Resource: "reserved-instances",
				},
				{
					Arn:           "arn:${Partition}:iam::${Account}:role/${RoleNameWithPath}",
					ConditionKeys: []string{},
					Resource:      "role",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:route-table/${RouteTableId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:Vpc",
					},
					Resource: "route-table",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:security-group/${SecurityGroupId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:Vpc",
					},
					Resource: "security-group",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:security-group-rule/${SecurityGroupRuleId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "security-group-rule",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}::snapshot/${SnapshotId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:OutpostArn",
						"ec2:Owner",
						"ec2:ParentVolume",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:SnapshotTime",
						"ec2:SourceOutpostArn",
						"ec2:VolumeSize",
					},
					Resource: "snapshot",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:spot-fleet-request/${SpotFleetRequestId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "spot-fleet-request",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:spot-instances-request/${SpotInstanceRequestId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "spot-instances-request",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:subnet/${SubnetId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:AvailabilityZone",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:Vpc",
					},
					Resource: "subnet",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:traffic-mirror-filter/${TrafficMirrorFilterId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "traffic-mirror-filter",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:traffic-mirror-filter-rule/${TrafficMirrorFilterRuleId}",
					ConditionKeys: []string{
						"ec2:Region",
					},
					Resource: "traffic-mirror-filter-rule",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:traffic-mirror-session/${TrafficMirrorSessionId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "traffic-mirror-session",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:traffic-mirror-target/${TrafficMirrorTargetId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "traffic-mirror-target",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:transit-gateway-attachment/${TransitGatewayAttachmentId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "transit-gateway-attachment",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:transit-gateway-connect-peer/${TransitGatewayConnectPeerId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "transit-gateway-connect-peer",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:transit-gateway/${TransitGatewayId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "transit-gateway",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:transit-gateway-multicast-domain/${TransitGatewayMulticastDomainId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "transit-gateway-multicast-domain",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:transit-gateway-route-table/${TransitGatewayRouteTableId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "transit-gateway-route-table",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:volume/${VolumeId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:AvailabilityZone",
						"ec2:Encrypted",
						"ec2:ParentSnapshot",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:VolumeIops",
						"ec2:VolumeSize",
						"ec2:VolumeThroughput",
						"ec2:VolumeType",
					},
					Resource: "volume",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:vpc-endpoint/${VpcEndpointId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "vpc-endpoint",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:vpc-endpoint-service/${VpcEndpointServiceId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:VpceServicePrivateDnsName",
					},
					Resource: "vpc-endpoint-service",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:vpc-flow-log/${VpcFlowLogId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "vpc-flow-log",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:vpc/${VpcId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
						"ec2:Tenancy",
					},
					Resource: "vpc",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:vpc-peering-connection/${VpcPeeringConnectionId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:AccepterVpc",
						"ec2:Region",
						"ec2:RequesterVpc",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "vpc-peering-connection",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:vpn-connection/${VpnConnectionId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:AuthenticationType",
						"ec2:DPDTimeoutSeconds",
						"ec2:GatewayType",
						"ec2:IKEVersions",
						"ec2:InsideTunnelCidr",
						"ec2:Phase1DHGroupNumbers",
						"ec2:Phase1EncryptionAlgorithms",
						"ec2:Phase1IntegrityAlgorithms",
						"ec2:Phase1LifetimeSeconds",
						"ec2:Phase2DHGroupNumbers",
						"ec2:Phase2EncryptionAlgorithms",
						"ec2:Phase2IntegrityAlgorithms",
						"ec2:Phase2LifetimeSeconds",
						"ec2:PresharedKeys",
						"ec2:Region",
						"ec2:RekeyFuzzPercentage",
						"ec2:RekeyMarginTimeSeconds",
						"ec2:ResourceTag/${TagKey}",
						"ec2:RoutingType",
					},
					Resource: "vpn-connection",
				},
				{
					Arn: "arn:${Partition}:ec2:${Region}:${Account}:vpn-gateway/${VpnGatewayId}",
					ConditionKeys: []string{
						"aws:RequestTag/${TagKey}",
						"aws:ResourceTag/${TagKey}",
						"aws:TagKeys",
						"ec2:Region",
						"ec2:ResourceTag/${TagKey}",
					},
					Resource: "vpn-gateway",
				},
			},
			ServiceName: "Amazon EC2",
		},
		ParliamentService{
			Conditions: []ParliamentCondition{
				{
					Condition:   "aws:RequestTag/${TagKey}",
					Description: "Filters actions based on the tags that are passed in the request",
					Type:        "String",
				},
				{
					Condition:   "aws:RequestedRegion",
					Description: "Requested region for the multi region access point operation",
					Type:        "String",
				},
				{
					Condition:   "aws:ResourceTag/${TagKey}",
					Description: "Filters actions based on the tags associated with the resource",
					Type:        "String",
				},
				{
					Condition:   "aws:TagKeys",
					Description: "Filters actions based on the tag keys that are passed in the request",
					Type:        "ArrayOfString",
				},
				{
					Condition:   "s3:AccessPointNetworkOrigin",
					Description: "Filters access by the network origin (Internet or VPC)",
					Type:        "String",
				},
				{
					Condition:   "s3:DataAccessPointAccount",
					Description: "Filters access by the AWS Account ID that owns the access point",
					Type:        "String",
				},
				{
					Condition:   "s3:DataAccessPointArn",
					Description: "Filters access by an access point Amazon Resource Name (ARN)",
					Type:        "String",
				},
				{
					Condition:   "s3:ExistingJobOperation",
					Description: "Filters access to updating the job priority by operation",
					Type:        "String",
				},
				{
					Condition:   "s3:ExistingJobPriority",
					Description: "Filters access to cancelling existing jobs by priority range",
					Type:        "Numeric",
				},
				{
					Condition:   "s3:ExistingObjectTag/<key>",
					Description: "Filters access by existing object tag key and value",
					Type:        "String",
				},
				{
					Condition:   "s3:JobSuspendedCause",
					Description: "Filters access to cancelling suspended jobs by a specific job suspended cause (for example, AWAITING_CONFIRMATION)",
					Type:        "String",
				},
				{
					Condition:   "s3:LocationConstraint",
					Description: "Filters access by a specific Region",
					Type:        "String",
				},
				{
					Condition:   "s3:RequestJobOperation",
					Description: "Filters access to creating jobs by operation",
					Type:        "String",
				},
				{
					Condition:   "s3:RequestJobPriority",
					Description: "Filters access to creating new jobs by priority range",
					Type:        "Numeric",
				},
				{
					Condition:   "s3:RequestObjectTag/<key>",
					Description: "Filters access by the tag keys and values to be added to objects",
					Type:        "String",
				},
				{
					Condition:   "s3:RequestObjectTagKeys",
					Description: "Filters access by the tag keys to be added to objects",
					Type:        "ArrayOfString",
				},
				{
					Condition:   "s3:ResourceAccount",
					Description: "Filters access by the resource owner AWS account ID",
					Type:        "String",
				},
				{
					Condition:   "s3:TlsVersion",
					Description: "Filters access by the TLS version used by the client",
					Type:        "Numeric",
				},
				{
					Condition:   "s3:VersionId",
					Description: "Filters access by a specific object version",
					Type:        "String",
				},
				{
					Condition:   "s3:authType",
					Description: "Filters access by authentication method",
					Type:        "String",
				},
				{
					Condition:   "s3:delimiter",
					Description: "Filters access by delimiter parameter",
					Type:        "String",
				},
				{
					Condition:   "s3:locationconstraint",
					Description: "Filters access by a specific Region",
					Type:        "String",
				},
				{
					Condition:   "s3:max-keys",
					Description: "Filters access by maximum number of keys returned in a ListBucket request",
					Type:        "Numeric",
				},
				{
					Condition:   "s3:object-lock-legal-hold",
					Description: "Filters access by object legal hold status",
					Type:        "String",
				},
				{
					Condition:   "s3:object-lock-mode",
					Description: "Filters access by object retention mode (COMPLIANCE or GOVERNANCE)",
					Type:        "String",
				},
				{
					Condition:   "s3:object-lock-remaining-retention-days",
					Description: "Filters access by remaining object retention days",
					Type:        "String",
				},
				{
					Condition:   "s3:object-lock-retain-until-date",
					Description: "Filters access by object retain-until date",
					Type:        "String",
				},
				{
					Condition:   "s3:prefix",
					Description: "Filters access by key name prefix",
					Type:        "String",
				},
				{
					Condition:   "s3:signatureAge",
					Description: "Filters access by the age in milliseconds of the request signature",
					Type:        "Numeric",
				},
				{
					Condition:   "s3:signatureversion",
					Description: "Filters access by the version of AWS Signature used on the request",
					Type:        "String",
				},
				{
					Condition:   "s3:versionid",
					Description: "Filters access by a specific object version",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-acl",
					Description: "Filters access by canned ACL in the request's x-amz-acl header",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-content-sha256",
					Description: "Filters access to unsigned content in your bucket",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-copy-source",
					Description: "Filters access to requests with a specific bucket, prefix, or object as the copy source",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-grant-full-control",
					Description: "Filters access to requests with the x-amz-grant-full-control (full control) header",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-grant-read",
					Description: "Filters access to requests with the x-amz-grant-read (read access) header",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-grant-read-acp",
					Description: "Filters access to requests with the x-amz-grant-read-acp (read permissions for the ACL) header",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-grant-write",
					Description: "Filters access to requests with the x-amz-grant-write (write access) header",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-grant-write-acp",
					Description: "Filters access to requests with the x-amz-grant-write-acp (write permissions for the ACL) header",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-metadata-directive",
					Description: "Filters access by object metadata behavior (COPY or REPLACE) when objects are copied",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-server-side-encryption",
					Description: "Filters access by server-side encryption",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-server-side-encryption-aws-kms-key-id",
					Description: "Filters access by AWS KMS customer managed CMK for server-side encryption",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-server-side-encryption-customer-algorithm",
					Description: "Filters access by customer-provided algorithm (SSE-C) for server-side encryption",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-storage-class",
					Description: "Filters access by storage class",
					Type:        "String",
				},
				{
					Condition:   "s3:x-amz-website-redirect-location",
					Description: "Filters access by a specific website redirect location for buckets that are configured as static websites",
					Type:        "String",
				},
			},
			Prefix: "s3",
			Privileges: []ParliamentPrivilege{
				{
					AccessLevel: "Write",
					Description: "Grants permission to abort a multipart upload",
					Privilege:   "AbortMultipartUpload",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointArn",
								"s3:DataAccessPointAccount",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to allow circumvention of governance-mode object retention settings",
					Privilege:   "BypassGovernanceRetention",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:RequestObjectTag/<key>",
								"s3:RequestObjectTagKeys",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-acl",
								"s3:x-amz-content-sha256",
								"s3:x-amz-copy-source",
								"s3:x-amz-grant-full-control",
								"s3:x-amz-grant-read",
								"s3:x-amz-grant-read-acp",
								"s3:x-amz-grant-write",
								"s3:x-amz-grant-write-acp",
								"s3:x-amz-metadata-directive",
								"s3:x-amz-server-side-encryption",
								"s3:x-amz-server-side-encryption-aws-kms-key-id",
								"s3:x-amz-server-side-encryption-customer-algorithm",
								"s3:x-amz-storage-class",
								"s3:x-amz-website-redirect-location",
								"s3:object-lock-mode",
								"s3:object-lock-retain-until-date",
								"s3:object-lock-remaining-retention-days",
								"s3:object-lock-legal-hold",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new access point",
					Privilege:   "CreateAccessPoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "accesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:locationconstraint",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-acl",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an object lambda enabled accesspoint",
					Privilege:   "CreateAccessPointForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "objectlambdaaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new bucket",
					Privilege:   "CreateBucket",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:locationconstraint",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-acl",
								"s3:x-amz-content-sha256",
								"s3:x-amz-grant-full-control",
								"s3:x-amz-grant-read",
								"s3:x-amz-grant-read-acp",
								"s3:x-amz-grant-write",
								"s3:x-amz-grant-write-acp",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new Amazon S3 Batch Operations job",
					Privilege:   "CreateJob",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"s3:RequestJobPriority",
								"s3:RequestJobOperation",
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{
								"iam:PassRole",
							},
							ResourceType: "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new multi region access point",
					Privilege:   "CreateMultiRegionAccessPoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "multiregionaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"aws:RequestedRegion",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureversion",
								"s3:signatureAge",
								"s3:TlsVersion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the access point named in the URI",
					Privilege:   "DeleteAccessPoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "accesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointArn",
								"s3:DataAccessPointAccount",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the object lambda enabled access point named in the URI",
					Privilege:   "DeleteAccessPointForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "objectlambdaaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointArn",
								"s3:DataAccessPointAccount",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete the policy on a specified access point",
					Privilege:   "DeleteAccessPointPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "accesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointArn",
								"s3:DataAccessPointAccount",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete the policy on a specified object lambda enabled access point",
					Privilege:   "DeleteAccessPointPolicyForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "objectlambdaaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointArn",
								"s3:DataAccessPointAccount",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the bucket named in the URI",
					Privilege:   "DeleteBucket",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete ownership controls on a bucket",
					Privilege:   "DeleteBucketOwnershipControls",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete the policy on a specified bucket",
					Privilege:   "DeleteBucketPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove the website configuration for a bucket",
					Privilege:   "DeleteBucketWebsite",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove tags from an existing Amazon S3 Batch Operations job",
					Privilege:   "DeleteJobTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "job*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"s3:ExistingJobPriority",
								"s3:ExistingJobOperation",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the multi region access point named in the URI",
					Privilege:   "DeleteMultiRegionAccessPoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "multiregionaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"aws:RequestedRegion",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureversion",
								"s3:signatureAge",
								"s3:TlsVersion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove the null version of an object and insert a delete marker, which becomes the current version of the object",
					Privilege:   "DeleteObject",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to use the tagging subresource to remove the entire tag set from the specified object",
					Privilege:   "DeleteObjectTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove a specific version of an object",
					Privilege:   "DeleteObjectVersion",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:versionid",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove the entire tag set for a specific version of the object",
					Privilege:   "DeleteObjectVersionTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:versionid",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an existing Amazon S3 Storage Lens configuration",
					Privilege:   "DeleteStorageLensConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "storagelensconfiguration*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove tags from an existing Amazon S3 Storage Lens configuration",
					Privilege:   "DeleteStorageLensConfigurationTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "storagelensconfiguration*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the configuration parameters and status for a batch operations job",
					Privilege:   "DescribeJob",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "job*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the configurations for a multi region access point",
					Privilege:   "DescribeMultiRegionAccessPointOperation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "multiregionaccesspointrequestarn*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestedRegion",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureversion",
								"s3:signatureAge",
								"s3:TlsVersion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to uses the accelerate subresource to return the Transfer Acceleration state of a bucket, which is either Enabled or Suspended",
					Privilege:   "GetAccelerateConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return configuration information about the specified access point",
					Privilege:   "GetAccessPoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the configuration of the object lambda enabled access point",
					Privilege:   "GetAccessPointConfigurationForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "objectlambdaaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointArn",
								"s3:DataAccessPointAccount",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to create an object lambda enabled accesspoint",
					Privilege:   "GetAccessPointForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "objectlambdaaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to returns the access point policy associated with the specified access point",
					Privilege:   "GetAccessPointPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "accesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to returns the access point policy associated with the specified object lambda enabled access point",
					Privilege:   "GetAccessPointPolicyForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "objectlambdaaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the policy status for a specific access point policy",
					Privilege:   "GetAccessPointPolicyStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "accesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the policy status for a specific object lambda access point policy",
					Privilege:   "GetAccessPointPolicyStatusForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "objectlambdaaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the PublicAccessBlock configuration for an AWS account",
					Privilege:   "GetAccountPublicAccessBlock",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get an analytics configuration from an Amazon S3 bucket, identified by the analytics configuration ID",
					Privilege:   "GetAnalyticsConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to use the acl subresource to return the access control list (ACL) of an Amazon S3 bucket",
					Privilege:   "GetBucketAcl",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the CORS configuration information set for an Amazon S3 bucket",
					Privilege:   "GetBucketCORS",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the Region that an Amazon S3 bucket resides in",
					Privilege:   "GetBucketLocation",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the logging status of an Amazon S3 bucket and the permissions users have to view or modify that status",
					Privilege:   "GetBucketLogging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get the notification configuration of an Amazon S3 bucket",
					Privilege:   "GetBucketNotification",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get the Object Lock configuration of an Amazon S3 bucket",
					Privilege:   "GetBucketObjectLockConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:signatureversion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve ownership controls on a bucket",
					Privilege:   "GetBucketOwnershipControls",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the policy of the specified bucket",
					Privilege:   "GetBucketPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the policy status for a specific Amazon S3 bucket, which indicates whether the bucket is public",
					Privilege:   "GetBucketPolicyStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the PublicAccessBlock configuration for an Amazon S3 bucket",
					Privilege:   "GetBucketPublicAccessBlock",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the request payment configuration for an Amazon S3 bucket",
					Privilege:   "GetBucketRequestPayment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the tag set associated with an Amazon S3 bucket",
					Privilege:   "GetBucketTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the versioning state of an Amazon S3 bucket",
					Privilege:   "GetBucketVersioning",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the website configuration for an Amazon S3 bucket",
					Privilege:   "GetBucketWebsite",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the default encryption configuration an Amazon S3 bucket",
					Privilege:   "GetEncryptionConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get an or list all Amazon S3 Intelligent Tiering configuration in a S3 Bucket",
					Privilege:   "GetIntelligentTieringConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return an inventory configuration from an Amazon S3 bucket, identified by the inventory configuration ID",
					Privilege:   "GetInventoryConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the tag set of an existing Amazon S3 Batch Operations job",
					Privilege:   "GetJobTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "job*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the lifecycle configuration information set on an Amazon S3 bucket",
					Privilege:   "GetLifecycleConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get a metrics configuration from an Amazon S3 bucket",
					Privilege:   "GetMetricsConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return configuration information about the specified multi region access point",
					Privilege:   "GetMultiRegionAccessPoint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "multiregionaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"aws:RequestedRegion",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureversion",
								"s3:signatureAge",
								"s3:TlsVersion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to returns the access point policy associated with the specified multi region access point",
					Privilege:   "GetMultiRegionAccessPointPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "multiregionaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"aws:RequestedRegion",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureversion",
								"s3:signatureAge",
								"s3:TlsVersion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the policy status for a specific multi region access point policy",
					Privilege:   "GetMultiRegionAccessPointPolicyStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "multiregionaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"aws:RequestedRegion",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureversion",
								"s3:signatureAge",
								"s3:TlsVersion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve objects from Amazon S3",
					Privilege:   "GetObject",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the access control list (ACL) of an object",
					Privilege:   "GetObjectAcl",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get an object's current Legal Hold status",
					Privilege:   "GetObjectLegalHold",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the retention settings for an object",
					Privilege:   "GetObjectRetention",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the tag set of an object",
					Privilege:   "GetObjectTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return torrent files from an Amazon S3 bucket",
					Privilege:   "GetObjectTorrent",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve a specific version of an object",
					Privilege:   "GetObjectVersion",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:versionid",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the access control list (ACL) of a specific object version",
					Privilege:   "GetObjectVersionAcl",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:versionid",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to replicate both unencrypted objects and objects encrypted with SSE-S3 or SSE-KMS",
					Privilege:   "GetObjectVersionForReplication",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to return the tag set for a specific version of the object",
					Privilege:   "GetObjectVersionTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:versionid",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get Torrent files about a different version using the versionId subresource",
					Privilege:   "GetObjectVersionTorrent",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:versionid",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get the replication configuration information set on an Amazon S3 bucket",
					Privilege:   "GetReplicationConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get an Amazon S3 Storage Lens configuration",
					Privilege:   "GetStorageLensConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "storagelensconfiguration*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get the tag set of an existing Amazon S3 Storage Lens configuration",
					Privilege:   "GetStorageLensConfigurationTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "storagelensconfiguration*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get an Amazon S3 Storage Lens dashboard",
					Privilege:   "GetStorageLensDashboard",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "storagelensconfiguration*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list access points",
					Privilege:   "ListAccessPoints",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list object lambda enabled accesspoints",
					Privilege:   "ListAccessPointsForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list all buckets owned by the authenticated sender of the request",
					Privilege:   "ListAllMyBuckets",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list some or all of the objects in an Amazon S3 bucket (up to 1000)",
					Privilege:   "ListBucket",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:delimiter",
								"s3:max-keys",
								"s3:prefix",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list in-progress multipart uploads",
					Privilege:   "ListBucketMultipartUploads",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list metadata about all the versions of objects in an Amazon S3 bucket",
					Privilege:   "ListBucketVersions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:delimiter",
								"s3:max-keys",
								"s3:prefix",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list current jobs and jobs that have ended recently",
					Privilege:   "ListJobs",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list multi region access points",
					Privilege:   "ListMultiRegionAccessPoints",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"aws:RequestedRegion",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureversion",
								"s3:signatureAge",
								"s3:TlsVersion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the parts that have been uploaded for a specific multipart upload",
					Privilege:   "ListMultipartUploadParts",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list Amazon S3 Storage Lens configurations",
					Privilege:   "ListStorageLensConfigurations",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to change replica ownership",
					Privilege:   "ObjectOwnerOverrideToBucketOwner",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to use the accelerate subresource to set the Transfer Acceleration state of an existing S3 bucket",
					Privilege:   "PutAccelerateConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set the configuration of the object lambda enabled access point",
					Privilege:   "PutAccessPointConfigurationForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "objectlambdaaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointArn",
								"s3:DataAccessPointAccount",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to associate an access policy with a specified access point",
					Privilege:   "PutAccessPointPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "accesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to associate an access policy with a specified object lambda enabled access point",
					Privilege:   "PutAccessPointPolicyForObjectLambda",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "objectlambdaaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to create or modify the PublicAccessBlock configuration for an AWS account",
					Privilege:   "PutAccountPublicAccessBlock",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set an analytics configuration for the bucket, specified by the analytics configuration ID",
					Privilege:   "PutAnalyticsConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to set the permissions on an existing bucket using access control lists (ACLs)",
					Privilege:   "PutBucketAcl",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-acl",
								"s3:x-amz-content-sha256",
								"s3:x-amz-grant-full-control",
								"s3:x-amz-grant-read",
								"s3:x-amz-grant-read-acp",
								"s3:x-amz-grant-write",
								"s3:x-amz-grant-write-acp",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set the CORS configuration for an Amazon S3 bucket",
					Privilege:   "PutBucketCORS",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set the logging parameters for an Amazon S3 bucket",
					Privilege:   "PutBucketLogging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to receive notifications when certain events happen in an Amazon S3 bucket",
					Privilege:   "PutBucketNotification",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to put Object Lock configuration on a specific bucket",
					Privilege:   "PutBucketObjectLockConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:TlsVersion",
								"s3:signatureversion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add or replace ownership controls on a bucket",
					Privilege:   "PutBucketOwnershipControls",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to add or replace a bucket policy on a bucket",
					Privilege:   "PutBucketPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to create or modify the PublicAccessBlock configuration for a specific Amazon S3 bucket",
					Privilege:   "PutBucketPublicAccessBlock",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set the request payment configuration of a bucket",
					Privilege:   "PutBucketRequestPayment",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add a set of tags to an existing Amazon S3 bucket",
					Privilege:   "PutBucketTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set the versioning state of an existing Amazon S3 bucket",
					Privilege:   "PutBucketVersioning",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set the configuration of the website that is specified in the website subresource",
					Privilege:   "PutBucketWebsite",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set the encryption configuration for an Amazon S3 bucket",
					Privilege:   "PutEncryptionConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create new or update or delete an existing Amazon S3 Intelligent Tiering configuration",
					Privilege:   "PutIntelligentTieringConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add an inventory configuration to the bucket, identified by the inventory ID",
					Privilege:   "PutInventoryConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to replace tags on an existing Amazon S3 Batch Operations job",
					Privilege:   "PutJobTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "job*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"s3:ExistingJobPriority",
								"s3:ExistingJobOperation",
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new lifecycle configuration for the bucket or replace an existing lifecycle configuration",
					Privilege:   "PutLifecycleConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set or update a metrics configuration for the CloudWatch request metrics from an Amazon S3 bucket",
					Privilege:   "PutMetricsConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to associate an access policy with a specified multi region access point",
					Privilege:   "PutMultiRegionAccessPointPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "multiregionaccesspoint*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"aws:RequestedRegion",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureversion",
								"s3:signatureAge",
								"s3:TlsVersion",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add an object to a bucket",
					Privilege:   "PutObject",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:RequestObjectTag/<key>",
								"s3:RequestObjectTagKeys",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-acl",
								"s3:x-amz-content-sha256",
								"s3:x-amz-copy-source",
								"s3:x-amz-grant-full-control",
								"s3:x-amz-grant-read",
								"s3:x-amz-grant-read-acp",
								"s3:x-amz-grant-write",
								"s3:x-amz-grant-write-acp",
								"s3:x-amz-metadata-directive",
								"s3:x-amz-server-side-encryption",
								"s3:x-amz-server-side-encryption-aws-kms-key-id",
								"s3:x-amz-server-side-encryption-customer-algorithm",
								"s3:x-amz-storage-class",
								"s3:x-amz-website-redirect-location",
								"s3:object-lock-mode",
								"s3:object-lock-retain-until-date",
								"s3:object-lock-remaining-retention-days",
								"s3:object-lock-legal-hold",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to set the access control list (ACL) permissions for new or existing objects in an S3 bucket.",
					Privilege:   "PutObjectAcl",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-acl",
								"s3:x-amz-content-sha256",
								"s3:x-amz-grant-full-control",
								"s3:x-amz-grant-read",
								"s3:x-amz-grant-read-acp",
								"s3:x-amz-grant-write",
								"s3:x-amz-grant-write-acp",
								"s3:x-amz-storage-class",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to apply a Legal Hold configuration to the specified object",
					Privilege:   "PutObjectLegalHold",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"s3:object-lock-legal-hold",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to place an Object Retention configuration on an object",
					Privilege:   "PutObjectRetention",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"s3:object-lock-mode",
								"s3:object-lock-retain-until-date",
								"s3:object-lock-remaining-retention-days",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to set the supplied tag-set to an object that already exists in a bucket",
					Privilege:   "PutObjectTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:RequestObjectTag/<key>",
								"s3:RequestObjectTagKeys",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to use the acl subresource to set the access control list (ACL) permissions for an object that already exists in a bucket",
					Privilege:   "PutObjectVersionAcl",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:versionid",
								"s3:x-amz-acl",
								"s3:x-amz-content-sha256",
								"s3:x-amz-grant-full-control",
								"s3:x-amz-grant-read",
								"s3:x-amz-grant-read-acp",
								"s3:x-amz-grant-write",
								"s3:x-amz-grant-write-acp",
								"s3:x-amz-storage-class",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to set the supplied tag-set for a specific version of an object",
					Privilege:   "PutObjectVersionTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:ExistingObjectTag/<key>",
								"s3:RequestObjectTag/<key>",
								"s3:RequestObjectTagKeys",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:versionid",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new replication configuration or replace an existing one",
					Privilege:   "PutReplicationConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{},
							DependentActions: []string{
								"iam:PassRole",
							},
							ResourceType: "bucket*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create or update an Amazon S3 Storage Lens configuration",
					Privilege:   "PutStorageLensConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to put or replace tags on an existing Amazon S3 Storage Lens configuration",
					Privilege:   "PutStorageLensConfigurationTagging",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "storagelensconfiguration*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to replicate delete markers to the destination bucket",
					Privilege:   "ReplicateDelete",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to replicate objects and object tags to the destination bucket",
					Privilege:   "ReplicateObject",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"s3:x-amz-server-side-encryption",
								"s3:x-amz-server-side-encryption-aws-kms-key-id",
								"s3:x-amz-server-side-encryption-customer-algorithm",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to replicate object tags to the destination bucket",
					Privilege:   "ReplicateTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to restore an archived copy of an object back into Amazon S3",
					Privilege:   "RestoreObject",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "object*",
						},
						{
							ConditionKeys: []string{
								"s3:DataAccessPointAccount",
								"s3:DataAccessPointArn",
								"s3:AccessPointNetworkOrigin",
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the priority of an existing job",
					Privilege:   "UpdateJobPriority",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "job*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"s3:RequestJobPriority",
								"s3:ExistingJobPriority",
								"s3:ExistingJobOperation",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the status for the specified job",
					Privilege:   "UpdateJobStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "job*",
						},
						{
							ConditionKeys: []string{
								"s3:authType",
								"s3:ResourceAccount",
								"s3:signatureAge",
								"s3:signatureversion",
								"s3:TlsVersion",
								"s3:x-amz-content-sha256",
								"s3:ExistingJobPriority",
								"s3:ExistingJobOperation",
								"s3:JobSuspendedCause",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
			},
			Resources: []ParliamentResource{
				{
					Arn:           "arn:${Partition}:s3:${Region}:${Account}:accesspoint/${AccessPointName}",
					ConditionKeys: []string{},
					Resource:      "accesspoint",
				},
				{
					Arn:           "arn:${Partition}:s3:::${BucketName}",
					ConditionKeys: []string{},
					Resource:      "bucket",
				},
				{
					Arn:           "arn:${Partition}:s3:::${BucketName}/${ObjectName}",
					ConditionKeys: []string{},
					Resource:      "object",
				},
				{
					Arn:           "arn:${Partition}:s3:${Region}:${Account}:job/${JobId}",
					ConditionKeys: []string{},
					Resource:      "job",
				},
				{
					Arn: "arn:${Partition}:s3:${Region}:${Account}:storage-lens/${ConfigId}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
					},
					Resource: "storagelensconfiguration",
				},
				{
					Arn:           "arn:${Partition}:s3-object-lambda:${Region}:${Account}:accesspoint/${AccessPointName}",
					ConditionKeys: []string{},
					Resource:      "objectlambdaaccesspoint",
				},
				{
					Arn:           "arn:${Partition}:s3::${Account}:accesspoint/${AccessPointAlias}",
					ConditionKeys: []string{},
					Resource:      "multiregionaccesspoint",
				},
				{
					Arn:           "arn:${Partition}:s3:us-west-2:${Account}:async-request/mrap/${Operation}/${Token}",
					ConditionKeys: []string{},
					Resource:      "multiregionaccesspointrequestarn",
				},
			},
			ServiceName: "Amazon S3",
		},
		ParliamentService{
			Conditions: []ParliamentCondition{
				{
					Condition:   "aws:RequestTag/${TagKey}",
					Description: "Filters create requests based on the allowed set of values for each of the tags",
					Type:        "String",
				},
				{
					Condition:   "aws:ResourceTag/${TagKey}",
					Description: "Filters actions based on tag-value associated with the resource",
					Type:        "String",
				},
				{
					Condition:   "aws:TagKeys",
					Description: "Filters create requests based on the presence of mandatory tags in the request",
					Type:        "String",
				},
				{
					Condition:   "ecr-public:ResourceTag/${TagKey}",
					Description: "Filters actions based on tag-value associated with the resource",
					Type:        "String",
				},
			},
			Prefix: "ecr-public",
			Privileges: []ParliamentPrivilege{
				{
					AccessLevel: "Read",
					Description: "Grants permission to check the availability of multiple image layers in a specified registry and repository",
					Privilege:   "BatchCheckLayerAvailability",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a list of specified images within a specified repository",
					Privilege:   "BatchDeleteImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to inform Amazon ECR that the image layer upload for a specified registry, repository name, and upload ID, has completed",
					Privilege:   "CompleteLayerUpload",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an image repository",
					Privilege:   "CreateRepository",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an existing image repository",
					Privilege:   "DeleteRepository",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the repository policy from a specified repository",
					Privilege:   "DeleteRepositoryPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe all the image tags for a given repository",
					Privilege:   "DescribeImageTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get metadata about the images in a repository, including image size, image tags, and creation date",
					Privilege:   "DescribeImages",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to retrieve the catalog data associated with a registry",
					Privilege:   "DescribeRegistries",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "registry*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to describe image repositories in a registry",
					Privilege:   "DescribeRepositories",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve a token that is valid for a specified registry for 12 hours",
					Privilege:   "GetAuthorizationToken",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the catalog data associated with a registry",
					Privilege:   "GetRegistryCatalogData",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "registry*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the catalog data associated with a repository",
					Privilege:   "GetRepositoryCatalogData",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the repository policy for a specified repository",
					Privilege:   "GetRepositoryPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to notify Amazon ECR that you intend to upload an image layer",
					Privilege:   "InitiateLayerUpload",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to list the tags for an Amazon ECR resource",
					Privilege:   "ListTagsForResource",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create or update the image manifest associated with an image",
					Privilege:   "PutImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create and update the catalog data associated with a registry",
					Privilege:   "PutRegistryCatalogData",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "registry*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the catalog data associated with a repository",
					Privilege:   "PutRepositoryCatalogData",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to apply a repository policy on a specified repository to control access permissions",
					Privilege:   "SetRepositoryPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to tag an Amazon ECR resource",
					Privilege:   "TagResource",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to untag an Amazon ECR resource",
					Privilege:   "UntagResource",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to upload an image layer part to Amazon ECR Public",
					Privilege:   "UploadLayerPart",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
			},
			Resources: []ParliamentResource{
				{
					Arn: "arn:${Partition}:ecr-public::${Account}:repository/${RepositoryName}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
						"ecr-public:ResourceTag/${TagKey}",
					},
					Resource: "repository",
				},
				{
					Arn:           "arn:${Partition}:ecr-public::${Account}:registry/${RegistryId}",
					ConditionKeys: []string{},
					Resource:      "registry",
				},
			},
			ServiceName: "Amazon Elastic Container Registry Public",
		},
		ParliamentService{
			Conditions: []ParliamentCondition{
				{
					Condition:   "aws:RequestTag/${TagKey}",
					Description: "Filters access based on the tags that are passed in the request",
					Type:        "String",
				},
				{
					Condition:   "aws:ResourceTag/${TagKey}",
					Description: "Filters access based on the tags associated with the resource",
					Type:        "String",
				},
				{
					Condition:   "aws:TagKeys",
					Description: "Filters access based on the tag keys that are passed in the request",
					Type:        "String",
				},
				{
					Condition:   "iam:AWSServiceName",
					Description: "Filters access by the AWS service to which this role is attached",
					Type:        "String",
				},
				{
					Condition:   "iam:AssociatedResourceArn",
					Description: "Filters by the resource that the role will be used on behalf of",
					Type:        "ARN",
				},
				{
					Condition:   "iam:OrganizationsPolicyId",
					Description: "Filters access by the ID of an AWS Organizations policy",
					Type:        "String",
				},
				{
					Condition:   "iam:PassedToService",
					Description: "Filters access by the AWS service to which this role is passed",
					Type:        "String",
				},
				{
					Condition:   "iam:PermissionsBoundary",
					Description: "Filters access if the specified policy is set as the permissions boundary on the IAM entity (user or role)",
					Type:        "String",
				},
				{
					Condition:   "iam:PolicyARN",
					Description: "Filters access by the ARN of an IAM policy",
					Type:        "ARN",
				},
				{
					Condition:   "iam:ResourceTag/${TagKey}",
					Description: "Filters access by the tags attached to an IAM entity (user or role)",
					Type:        "String",
				},
			},
			Prefix: "iam",
			Privileges: []ParliamentPrivilege{
				{
					AccessLevel: "Write",
					Description: "Grants permission to add a new client ID (audience) to the list of registered IDs for the specified IAM OpenID Connect (OIDC) provider resource",
					Privilege:   "AddClientIDToOpenIDConnectProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "oidc-provider*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add an IAM role to the specified instance profile",
					Privilege:   "AddRoleToInstanceProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{},
							DependentActions: []string{
								"iam:PassRole",
							},
							ResourceType: "instance-profile*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to add an IAM user to the specified IAM group",
					Privilege:   "AddUserToGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to attach a managed policy to the specified IAM group",
					Privilege:   "AttachGroupPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
						{
							ConditionKeys: []string{
								"iam:PolicyARN",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to attach a managed policy to the specified IAM role",
					Privilege:   "AttachRolePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"iam:PolicyARN",
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to attach a managed policy to the specified IAM user",
					Privilege:   "AttachUserPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
						{
							ConditionKeys: []string{
								"iam:PolicyARN",
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission for an IAM user to to change their own password",
					Privilege:   "ChangePassword",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create access key and secret access key for the specified IAM user",
					Privilege:   "CreateAccessKey",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an alias for your AWS account",
					Privilege:   "CreateAccountAlias",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new group",
					Privilege:   "CreateGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new instance profile",
					Privilege:   "CreateInstanceProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "instance-profile*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a password for the specified IAM user",
					Privilege:   "CreateLoginProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an IAM resource that describes an identity provider (IdP) that supports OpenID Connect (OIDC)",
					Privilege:   "CreateOpenIDConnectProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "oidc-provider*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to create a new managed policy",
					Privilege:   "CreatePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to create a new version of the specified managed policy",
					Privilege:   "CreatePolicyVersion",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new role",
					Privilege:   "CreateRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an IAM resource that describes an identity provider (IdP) that supports SAML 2.0",
					Privilege:   "CreateSAMLProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "saml-provider*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an IAM role that allows an AWS service to perform actions on your behalf",
					Privilege:   "CreateServiceLinkedRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"iam:AWSServiceName",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new service-specific credential for an IAM user",
					Privilege:   "CreateServiceSpecificCredential",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new IAM user",
					Privilege:   "CreateUser",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create a new virtual MFA device",
					Privilege:   "CreateVirtualMFADevice",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "mfa*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to deactivate the specified MFA device and remove its association with the IAM user for which it was originally enabled",
					Privilege:   "DeactivateMFADevice",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the access key pair that is associated with the specified IAM user",
					Privilege:   "DeleteAccessKey",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified AWS account alias",
					Privilege:   "DeleteAccountAlias",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete the password policy for the AWS account",
					Privilege:   "DeleteAccountPasswordPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified IAM group",
					Privilege:   "DeleteGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete the specified inline policy from its group",
					Privilege:   "DeleteGroupPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified instance profile",
					Privilege:   "DeleteInstanceProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "instance-profile*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the password for the specified IAM user",
					Privilege:   "DeleteLoginProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an OpenID Connect identity provider (IdP) resource object in IAM",
					Privilege:   "DeleteOpenIDConnectProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "oidc-provider*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete the specified managed policy and remove it from any IAM entities (users, groups, or roles) to which it is attached",
					Privilege:   "DeletePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete a version from the specified managed policy",
					Privilege:   "DeletePolicyVersion",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified role",
					Privilege:   "DeleteRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to remove the permissions boundary from a role",
					Privilege:   "DeleteRolePermissionsBoundary",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete the specified inline policy from the specified role",
					Privilege:   "DeleteRolePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a SAML provider resource in IAM",
					Privilege:   "DeleteSAMLProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "saml-provider*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified SSH public key",
					Privilege:   "DeleteSSHPublicKey",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified server certificate",
					Privilege:   "DeleteServerCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "server-certificate*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an IAM role that is linked to a specific AWS service, if the service is no longer using it",
					Privilege:   "DeleteServiceLinkedRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified service-specific credential for an IAM user",
					Privilege:   "DeleteServiceSpecificCredential",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a signing certificate that is associated with the specified IAM user",
					Privilege:   "DeleteSigningCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified IAM user",
					Privilege:   "DeleteUser",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to remove the permissions boundary from the specified IAM user",
					Privilege:   "DeleteUserPermissionsBoundary",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to delete the specified inline policy from an IAM user",
					Privilege:   "DeleteUserPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a virtual MFA device",
					Privilege:   "DeleteVirtualMFADevice",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "mfa",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "sms-mfa",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to detach a managed policy from the specified IAM group",
					Privilege:   "DetachGroupPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
						{
							ConditionKeys: []string{
								"iam:PolicyARN",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to detach a managed policy from the specified role",
					Privilege:   "DetachRolePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"iam:PolicyARN",
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to detach a managed policy from the specified IAM user",
					Privilege:   "DetachUserPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
						{
							ConditionKeys: []string{
								"iam:PolicyARN",
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to enable an MFA device and associate it with the specified IAM user",
					Privilege:   "EnableMFADevice",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to generate a credential report for the AWS account",
					Privilege:   "GenerateCredentialReport",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to generate an access report for an AWS Organizations entity",
					Privilege:   "GenerateOrganizationsAccessReport",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys: []string{},
							DependentActions: []string{
								"organizations:DescribePolicy",
								"organizations:ListChildren",
								"organizations:ListParents",
								"organizations:ListPoliciesForTarget",
								"organizations:ListRoots",
								"organizations:ListTargetsForPolicy",
							},
							ResourceType: "access-report*",
						},
						{
							ConditionKeys: []string{
								"iam:OrganizationsPolicyId",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to generate a service last accessed data report for an IAM resource",
					Privilege:   "GenerateServiceLastAccessedDetails",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about when the specified access key was last used",
					Privilege:   "GetAccessKeyLastUsed",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about all IAM users, groups, roles, and policies in your AWS account, including their relationships to one another",
					Privilege:   "GetAccountAuthorizationDetails",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the password policy for the AWS account",
					Privilege:   "GetAccountPasswordPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to retrieve information about IAM entity usage and IAM quotas in the AWS account",
					Privilege:   "GetAccountSummary",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve a list of all of the context keys that are referenced in the specified policy",
					Privilege:   "GetContextKeysForCustomPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve a list of all context keys that are referenced in all IAM policies that are attached to the specified IAM identity (user, group, or role)",
					Privilege:   "GetContextKeysForPrincipalPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve a credential report for the AWS account",
					Privilege:   "GetCredentialReport",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve a list of IAM users in the specified IAM group",
					Privilege:   "GetGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve an inline policy document that is embedded in the specified IAM group",
					Privilege:   "GetGroupPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about the specified instance profile, including the instance profile's path, GUID, ARN, and role",
					Privilege:   "GetInstanceProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "instance-profile*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to retrieve the user name and password creation date for the specified IAM user",
					Privilege:   "GetLoginProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about the specified OpenID Connect (OIDC) provider resource in IAM",
					Privilege:   "GetOpenIDConnectProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "oidc-provider*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve an AWS Organizations access report",
					Privilege:   "GetOrganizationsAccessReport",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about the specified managed policy, including the policy's default version and the total number of identities to which the policy is attached",
					Privilege:   "GetPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about a version of the specified managed policy, including the policy document",
					Privilege:   "GetPolicyVersion",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about the specified role, including the role's path, GUID, ARN, and the role's trust policy",
					Privilege:   "GetRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve an inline policy document that is embedded with the specified IAM role",
					Privilege:   "GetRolePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the SAML provider metadocument that was uploaded when the IAM SAML provider resource was created or updated",
					Privilege:   "GetSAMLProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "saml-provider*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the specified SSH public key, including metadata about the key",
					Privilege:   "GetSSHPublicKey",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about the specified server certificate stored in IAM",
					Privilege:   "GetServerCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "server-certificate*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about the service last accessed data report",
					Privilege:   "GetServiceLastAccessedDetails",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about the entities from the service last accessed data report",
					Privilege:   "GetServiceLastAccessedDetailsWithEntities",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve an IAM service-linked role deletion status",
					Privilege:   "GetServiceLinkedRoleDeletionStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve information about the specified IAM user, including the user's creation date, path, unique ID, and ARN",
					Privilege:   "GetUser",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve an inline policy document that is embedded in the specified IAM user",
					Privilege:   "GetUserPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list information about the access key IDs that are associated with the specified IAM user",
					Privilege:   "ListAccessKeys",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the account alias that is associated with the AWS account",
					Privilege:   "ListAccountAliases",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list all managed policies that are attached to the specified IAM group",
					Privilege:   "ListAttachedGroupPolicies",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list all managed policies that are attached to the specified IAM role",
					Privilege:   "ListAttachedRolePolicies",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list all managed policies that are attached to the specified IAM user",
					Privilege:   "ListAttachedUserPolicies",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list all IAM identities to which the specified managed policy is attached",
					Privilege:   "ListEntitiesForPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the names of the inline policies that are embedded in the specified IAM group",
					Privilege:   "ListGroupPolicies",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the IAM groups that have the specified path prefix",
					Privilege:   "ListGroups",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the IAM groups that the specified IAM user belongs to",
					Privilege:   "ListGroupsForUser",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the tags that are attached to the specified instance profile",
					Privilege:   "ListInstanceProfileTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "instance-profile*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the instance profiles that have the specified path prefix",
					Privilege:   "ListInstanceProfiles",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "instance-profile*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the instance profiles that have the specified associated IAM role",
					Privilege:   "ListInstanceProfilesForRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the tags that are attached to the specified virtual mfa device",
					Privilege:   "ListMFADeviceTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "mfa*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the MFA devices for an IAM user",
					Privilege:   "ListMFADevices",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the tags that are attached to the specified OpenID Connect provider",
					Privilege:   "ListOpenIDConnectProviderTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "oidc-provider*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list information about the IAM OpenID Connect (OIDC) provider resource objects that are defined in the AWS account",
					Privilege:   "ListOpenIDConnectProviders",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list all managed policies",
					Privilege:   "ListPolicies",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list information about the policies that grant an entity access to a specific service",
					Privilege:   "ListPoliciesGrantingServiceAccess",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the tags that are attached to the specified managed policy",
					Privilege:   "ListPolicyTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list information about the versions of the specified managed policy, including the version that is currently set as the policy's default version",
					Privilege:   "ListPolicyVersions",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the names of the inline policies that are embedded in the specified IAM role",
					Privilege:   "ListRolePolicies",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the tags that are attached to the specified IAM role",
					Privilege:   "ListRoleTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the IAM roles that have the specified path prefix",
					Privilege:   "ListRoles",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the tags that are attached to the specified SAML provider",
					Privilege:   "ListSAMLProviderTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "saml-provider*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the SAML provider resources in IAM",
					Privilege:   "ListSAMLProviders",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list information about the SSH public keys that are associated with the specified IAM user",
					Privilege:   "ListSSHPublicKeys",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the tags that are attached to the specified server certificate",
					Privilege:   "ListServerCertificateTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "server-certificate*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the server certificates that have the specified path prefix",
					Privilege:   "ListServerCertificates",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the service-specific credentials that are associated with the specified IAM user",
					Privilege:   "ListServiceSpecificCredentials",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list information about the signing certificates that are associated with the specified IAM user",
					Privilege:   "ListSigningCertificates",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the names of the inline policies that are embedded in the specified IAM user",
					Privilege:   "ListUserPolicies",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the tags that are attached to the specified IAM user",
					Privilege:   "ListUserTags",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list the IAM users that have the specified path prefix",
					Privilege:   "ListUsers",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list virtual MFA devices by assignment status",
					Privilege:   "ListVirtualMFADevices",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to pass a role to a service",
					Privilege:   "PassRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"iam:AssociatedResourceArn",
								"iam:PassedToService",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to create or update an inline policy document that is embedded in the specified IAM group",
					Privilege:   "PutGroupPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to set a managed policy as a permissions boundary for a role",
					Privilege:   "PutRolePermissionsBoundary",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to create or update an inline policy document that is embedded in the specified IAM role",
					Privilege:   "PutRolePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to set a managed policy as a permissions boundary for an IAM user",
					Privilege:   "PutUserPermissionsBoundary",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to create or update an inline policy document that is embedded in the specified IAM user",
					Privilege:   "PutUserPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
						{
							ConditionKeys: []string{
								"iam:PermissionsBoundary",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove the client ID (audience) from the list of client IDs in the specified IAM OpenID Connect (OIDC) provider resource",
					Privilege:   "RemoveClientIDFromOpenIDConnectProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "oidc-provider*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove an IAM role from the specified EC2 instance profile",
					Privilege:   "RemoveRoleFromInstanceProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "instance-profile*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to remove an IAM user from the specified group",
					Privilege:   "RemoveUserFromGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to reset the password for an existing service-specific credential for an IAM user",
					Privilege:   "ResetServiceSpecificCredential",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to synchronize the specified MFA device with its IAM entity (user or role)",
					Privilege:   "ResyncMFADevice",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to set the version of the specified policy as the policy's default version",
					Privilege:   "SetDefaultPolicyVersion",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set the STS global endpoint token version",
					Privilege:   "SetSecurityTokenServicePreferences",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to simulate whether an identity-based policy or resource-based policy provides permissions for specific API operations and resources",
					Privilege:   "SimulateCustomPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to simulate whether an identity-based policy that is attached to a specified IAM entity (user or role) provides permissions for specific API operations and resources",
					Privilege:   "SimulatePrincipalPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add tags to an instance profile",
					Privilege:   "TagInstanceProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "instance-profile*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add tags to a virtual mfa device",
					Privilege:   "TagMFADevice",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "mfa*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add tags to an OpenID Connect provider",
					Privilege:   "TagOpenIDConnectProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "oidc-provider*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add tags to a managed policy",
					Privilege:   "TagPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add tags to an IAM role",
					Privilege:   "TagRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add tags to a SAML Provider",
					Privilege:   "TagSAMLProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "saml-provider*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add tags to a server certificate",
					Privilege:   "TagServerCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "server-certificate*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add tags to an IAM user",
					Privilege:   "TagUser",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove the specified tags from the instance profile",
					Privilege:   "UntagInstanceProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "instance-profile*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove the specified tags from the virtual mfa device",
					Privilege:   "UntagMFADevice",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "mfa*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove the specified tags from the OpenID Connect provider",
					Privilege:   "UntagOpenIDConnectProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "oidc-provider*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove the specified tags from the managed policy",
					Privilege:   "UntagPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "policy*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove the specified tags from the role",
					Privilege:   "UntagRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove the specified tags from the SAML Provider",
					Privilege:   "UntagSAMLProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "saml-provider*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove the specified tags from the server certificate",
					Privilege:   "UntagServerCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "server-certificate*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to remove the specified tags from the user",
					Privilege:   "UntagUser",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the status of the specified access key as Active or Inactive",
					Privilege:   "UpdateAccessKey",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the password policy settings for the AWS account",
					Privilege:   "UpdateAccountPasswordPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to update the policy that grants an IAM entity permission to assume a role",
					Privilege:   "UpdateAssumeRolePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the name or path of the specified IAM group",
					Privilege:   "UpdateGroup",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "group*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to change the password for the specified IAM user",
					Privilege:   "UpdateLoginProfile",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the entire list of server certificate thumbprints that are associated with an OpenID Connect (OIDC) provider resource",
					Privilege:   "UpdateOpenIDConnectProviderThumbprint",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "oidc-provider*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the description or maximum session duration setting of a role",
					Privilege:   "UpdateRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update only the description of a role",
					Privilege:   "UpdateRoleDescription",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the metadata document for an existing SAML provider resource",
					Privilege:   "UpdateSAMLProvider",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "saml-provider*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the status of an IAM user's SSH public key to active or inactive",
					Privilege:   "UpdateSSHPublicKey",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the name or the path of the specified server certificate stored in IAM",
					Privilege:   "UpdateServerCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "server-certificate*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the status of a service-specific credential to active or inactive for an IAM user",
					Privilege:   "UpdateServiceSpecificCredential",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the status of the specified user signing certificate to active or disabled",
					Privilege:   "UpdateSigningCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the name or the path of the specified IAM user",
					Privilege:   "UpdateUser",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to upload an SSH public key and associate it with the specified IAM user",
					Privilege:   "UploadSSHPublicKey",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to upload a server certificate entity for the AWS account",
					Privilege:   "UploadServerCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "server-certificate*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to upload an X.509 signing certificate and associate it with the specified IAM user",
					Privilege:   "UploadSigningCertificate",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user*",
						},
					},
				},
			},
			Resources: []ParliamentResource{
				{
					Arn:           "arn:${Partition}:iam::${Account}:access-report/${EntityPath}",
					ConditionKeys: []string{},
					Resource:      "access-report",
				},
				{
					Arn:           "arn:${Partition}:iam::${Account}:assumed-role/${RoleName}/${RoleSessionName}",
					ConditionKeys: []string{},
					Resource:      "assumed-role",
				},
				{
					Arn:           "arn:${Partition}:iam::${Account}:federated-user/${UserName}",
					ConditionKeys: []string{},
					Resource:      "federated-user",
				},
				{
					Arn:           "arn:${Partition}:iam::${Account}:group/${GroupNameWithPath}",
					ConditionKeys: []string{},
					Resource:      "group",
				},
				{
					Arn: "arn:${Partition}:iam::${Account}:instance-profile/${InstanceProfileNameWithPath}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
					},
					Resource: "instance-profile",
				},
				{
					Arn: "arn:${Partition}:iam::${Account}:mfa/${MfaTokenIdWithPath}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
					},
					Resource: "mfa",
				},
				{
					Arn: "arn:${Partition}:iam::${Account}:oidc-provider/${OidcProviderName}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
					},
					Resource: "oidc-provider",
				},
				{
					Arn: "arn:${Partition}:iam::${Account}:policy/${PolicyNameWithPath}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
					},
					Resource: "policy",
				},
				{
					Arn: "arn:${Partition}:iam::${Account}:role/${RoleNameWithPath}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
						"iam:ResourceTag/${TagKey}",
					},
					Resource: "role",
				},
				{
					Arn: "arn:${Partition}:iam::${Account}:saml-provider/${SamlProviderName}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
					},
					Resource: "saml-provider",
				},
				{
					Arn: "arn:${Partition}:iam::${Account}:server-certificate/${CertificateNameWithPath}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
					},
					Resource: "server-certificate",
				},
				{
					Arn:           "arn:${Partition}:iam::${Account}:sms-mfa/${MfaTokenIdWithPath}",
					ConditionKeys: []string{},
					Resource:      "sms-mfa",
				},
				{
					Arn: "arn:${Partition}:iam::${Account}:user/${UserNameWithPath}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
						"iam:ResourceTag/${TagKey}",
					},
					Resource: "user",
				},
			},
			ServiceName: "Identity And Access Management",
		},
		ParliamentService{
			Conditions: []ParliamentCondition{
				{
					Condition:   "aws:RequestTag/${TagKey}",
					Description: "Filters create requests based on the allowed set of values for each of the tags",
					Type:        "String",
				},
				{
					Condition:   "aws:ResourceTag/${TagKey}",
					Description: "Filters actions based on tag-value associated with the resource",
					Type:        "String",
				},
				{
					Condition:   "aws:TagKeys",
					Description: "Filters create requests based on the presence of mandatory tags in the request",
					Type:        "String",
				},
				{
					Condition:   "ecr:ResourceTag/${TagKey}",
					Description: "Filters actions based on tag-value associated with the resource",
					Type:        "String",
				},
			},
			Prefix: "ecr",
			Privileges: []ParliamentPrivilege{
				{
					AccessLevel: "Read",
					Description: "Grants permission to check the availability of multiple image layers in a specified registry and repository",
					Privilege:   "BatchCheckLayerAvailability",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete a list of specified images within a specified repository",
					Privilege:   "BatchDeleteImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to get detailed information for specified images within a specified repository",
					Privilege:   "BatchGetImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to inform Amazon ECR that the image layer upload for a specified registry, repository name, and upload ID, has completed",
					Privilege:   "CompleteLayerUpload",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create an image repository",
					Privilege:   "CreateRepository",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the specified lifecycle policy",
					Privilege:   "DeleteLifecyclePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the registry policy",
					Privilege:   "DeleteRegistryPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete an existing image repository",
					Privilege:   "DeleteRepository",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to delete the repository policy from a specified repository",
					Privilege:   "DeleteRepositoryPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve replication status about an image in a registry, including failure reason if replication fails",
					Privilege:   "DescribeImageReplicationStatus",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe the image scan findings for the specified image",
					Privilege:   "DescribeImageScanFindings",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to get metadata about the images in a repository, including image size, image tags, and creation date",
					Privilege:   "DescribeImages",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe the registry settings",
					Privilege:   "DescribeRegistry",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to describe image repositories in a registry",
					Privilege:   "DescribeRepositories",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve a token that is valid for a specified registry for 12 hours",
					Privilege:   "GetAuthorizationToken",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the download URL corresponding to an image layer",
					Privilege:   "GetDownloadUrlForLayer",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the specified lifecycle policy",
					Privilege:   "GetLifecyclePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the results of the specified lifecycle policy preview request",
					Privilege:   "GetLifecyclePolicyPreview",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the registry policy",
					Privilege:   "GetRegistryPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to retrieve the repository policy for a specified repository",
					Privilege:   "GetRepositoryPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to notify Amazon ECR that you intend to upload an image layer",
					Privilege:   "InitiateLayerUpload",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "List",
					Description: "Grants permission to list all the image IDs for a given repository",
					Privilege:   "ListImages",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Grants permission to list the tags for an Amazon ECR resource",
					Privilege:   "ListTagsForResource",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create or update the image manifest associated with an image",
					Privilege:   "PutImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the image scanning configuration for a repository",
					Privilege:   "PutImageScanningConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the image tag mutability settings for a repository",
					Privilege:   "PutImageTagMutability",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to create or update a lifecycle policy",
					Privilege:   "PutLifecyclePolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the registry policy",
					Privilege:   "PutRegistryPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to update the replication configuration for the registry",
					Privilege:   "PutReplicationConfiguration",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to replicate images to the destination registry",
					Privilege:   "ReplicateImage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Permissions management",
					Description: "Grants permission to apply a repository policy on a specified repository to control access permissions",
					Privilege:   "SetRepositoryPolicy",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to start an image scan",
					Privilege:   "StartImageScan",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to start a preview of the specified lifecycle policy",
					Privilege:   "StartLifecyclePolicyPreview",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to tag an Amazon ECR resource",
					Privilege:   "TagResource",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
						{
							ConditionKeys: []string{
								"aws:RequestTag/${TagKey}",
								"aws:TagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to untag an Amazon ECR resource",
					Privilege:   "UntagResource",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to upload an image layer part to Amazon ECR",
					Privilege:   "UploadLayerPart",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "repository*",
						},
					},
				},
			},
			Resources: []ParliamentResource{
				{
					Arn: "arn:${Partition}:ecr:${Region}:${Account}:repository/${RepositoryName}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
						"ecr:ResourceTag/${TagKey}",
					},
					Resource: "repository",
				},
			},
			ServiceName: "Amazon Elastic Container Registry",
		},
		ParliamentService{
			Conditions: []ParliamentCondition{
				{
					Condition:   "accounts.google.com:aud",
					Description: "Filters actions based on the Google application ID",
					Type:        "String",
				},
				{
					Condition:   "accounts.google.com:oaud",
					Description: "Filters actions based on the Google audience",
					Type:        "String",
				},
				{
					Condition:   "accounts.google.com:sub",
					Description: "Filters actions based on the subject of the claim (the Google user ID)",
					Type:        "String",
				},
				{
					Condition:   "aws:FederatedProvider",
					Description: "Filters actions based on the IdP that was used to authenticate the user",
					Type:        "String",
				},
				{
					Condition:   "aws:PrincipalTag/${TagKey}",
					Description: "Filters actions based on the tag associated with the principal that is making the request",
					Type:        "String",
				},
				{
					Condition:   "aws:RequestTag/${TagKey}",
					Description: "Filters actions based on the tags that are passed in the request",
					Type:        "String",
				},
				{
					Condition:   "aws:ResourceTag/${TagKey}",
					Description: "Filters actions based on the tags associated with the resource",
					Type:        "String",
				},
				{
					Condition:   "aws:SourceIdentity",
					Description: "Filters actions based on the source identity that is set on the caller",
					Type:        "String",
				},
				{
					Condition:   "aws:TagKeys",
					Description: "Filters actions based on the tag keys that are passed in the request",
					Type:        "String",
				},
				{
					Condition:   "cognito-identity.amazonaws.com:amr",
					Description: "Filters actions based on the login information for Amazon Cognito",
					Type:        "String",
				},
				{
					Condition:   "cognito-identity.amazonaws.com:aud",
					Description: "Filters actions based on the Amazon Cognito identity pool ID",
					Type:        "String",
				},
				{
					Condition:   "cognito-identity.amazonaws.com:sub",
					Description: "Filters actions based on the subject of the claim (the Amazon Cognito user ID)",
					Type:        "String",
				},
				{
					Condition:   "graph.facebook.com:app_id",
					Description: "Filters actions based on the Facebook application ID",
					Type:        "String",
				},
				{
					Condition:   "graph.facebook.com:id",
					Description: "Filters actions based on the Facebook user ID",
					Type:        "String",
				},
				{
					Condition:   "iam:ResourceTag/${TagKey}",
					Description: "Filters actions based on the tags that are attached to the role that is being assumed",
					Type:        "String",
				},
				{
					Condition:   "saml:aud",
					Description: "Filters actions based on the endpoint URL to which SAML assertions are presented",
					Type:        "String",
				},
				{
					Condition:   "saml:cn",
					Description: "Filters actions based on the eduOrg attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:commonName",
					Description: "Filters actions based on the commonName attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:doc",
					Description: "Filters actions based on the principal that was used to assume the role",
					Type:        "String",
				},
				{
					Condition:   "saml:eduorghomepageuri",
					Description: "Filters actions based on the eduOrg attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:eduorgidentityauthnpolicyuri",
					Description: "Filters actions based on the eduOrg attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:eduorglegalname",
					Description: "Filters actions based on the eduOrg attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:eduorgsuperioruri",
					Description: "Filters actions based on the eduOrg attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:eduorgwhitepagesuri",
					Description: "Filters actions based on the eduOrg attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonaffiliation",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonassurance",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonentitlement",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonnickname",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonorgdn",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonorgunitdn",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonprimaryaffiliation",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonprimaryorgunitdn",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonprincipalname",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersonscopedaffiliation",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:edupersontargetedid",
					Description: "Filters actions based on the eduPerson attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:givenName",
					Description: "Filters actions based on the givenName attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:iss",
					Description: "Filters actions based on the issuer, which is represented by a URN",
					Type:        "String",
				},
				{
					Condition:   "saml:mail",
					Description: "Filters actions based on the mail attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:name",
					Description: "Filters actions based on the name attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:namequalifier",
					Description: "Filters actions based on the hash value of the issuer, account ID, and friendly name",
					Type:        "String",
				},
				{
					Condition:   "saml:organizationStatus",
					Description: "Filters actions based on the organizationStatus attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:primaryGroupSID",
					Description: "Filters actions based on the primaryGroupSID attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:sub",
					Description: "Filters actions based on the subject of the claim (the SAML user ID)",
					Type:        "String",
				},
				{
					Condition:   "saml:sub_type",
					Description: "Filters actions based on the value persistent, transient, or the full Format URI",
					Type:        "String",
				},
				{
					Condition:   "saml:surname",
					Description: "Filters actions based on the surname attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:uid",
					Description: "Filters actions based on the uid attribute",
					Type:        "String",
				},
				{
					Condition:   "saml:x500UniqueIdentifier",
					Description: "Filters actions based on the uid attribute",
					Type:        "String",
				},
				{
					Condition:   "sts:ExternalId",
					Description: "Filters actions based on the unique identifier required when you assume a role in another account",
					Type:        "String",
				},
				{
					Condition:   "sts:RoleSessionName",
					Description: "Filters actions based on the role session name required when you assume a role",
					Type:        "String",
				},
				{
					Condition:   "sts:SourceIdentity",
					Description: "Filters actions based on the source identity that is passed in the request",
					Type:        "String",
				},
				{
					Condition:   "sts:TransitiveTagKeys",
					Description: "Filters actions based on the transitive tag keys that are passed in the request",
					Type:        "String",
				},
				{
					Condition:   "www.amazon.com:app_id",
					Description: "Filters actions based on the Login with Amazon application ID",
					Type:        "String",
				},
				{
					Condition:   "www.amazon.com:user_id",
					Description: "Filters actions based on the Login with Amazon user ID",
					Type:        "String",
				},
			},
			Prefix: "sts",
			Privileges: []ParliamentPrivilege{
				{
					AccessLevel: "Write",
					Description: "Returns a set of temporary security credentials that you can use to access AWS resources that you might not normally have access to",
					Privilege:   "AssumeRole",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:PrincipalTag/${TagKey}",
								"aws:RequestTag/${TagKey}",
								"sts:TransitiveTagKeys",
								"sts:ExternalId",
								"sts:RoleSessionName",
								"iam:ResourceTag/${TagKey}",
								"sts:SourceIdentity",
								"aws:SourceIdentity",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Returns a set of temporary security credentials for users who have been authenticated via a SAML authentication response",
					Privilege:   "AssumeRoleWithSAML",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"saml:namequalifier",
								"saml:sub",
								"saml:sub_type",
								"saml:aud",
								"saml:iss",
								"saml:doc",
								"saml:cn",
								"saml:commonName",
								"saml:eduorghomepageuri",
								"saml:eduorgidentityauthnpolicyuri",
								"saml:eduorglegalname",
								"saml:eduorgsuperioruri",
								"saml:eduorgwhitepagesuri",
								"saml:edupersonaffiliation",
								"saml:edupersonassurance",
								"saml:edupersonentitlement",
								"saml:edupersonnickname",
								"saml:edupersonorgdn",
								"saml:edupersonorgunitdn",
								"saml:edupersonprimaryaffiliation",
								"saml:edupersonprimaryorgunitdn",
								"saml:edupersonprincipalname",
								"saml:edupersonscopedaffiliation",
								"saml:edupersontargetedid",
								"saml:givenName",
								"saml:mail",
								"saml:name",
								"saml:organizationStatus",
								"saml:primaryGroupSID",
								"saml:surname",
								"saml:uid",
								"saml:x500UniqueIdentifier",
								"aws:TagKeys",
								"aws:PrincipalTag/${TagKey}",
								"aws:RequestTag/${TagKey}",
								"sts:TransitiveTagKeys",
								"sts:SourceIdentity",
								"sts:RoleSessionName",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Returns a set of temporary security credentials for users who have been authenticated in a mobile or web application with a web identity provider",
					Privilege:   "AssumeRoleWithWebIdentity",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role*",
						},
						{
							ConditionKeys: []string{
								"cognito-identity.amazonaws.com:amr",
								"cognito-identity.amazonaws.com:aud",
								"cognito-identity.amazonaws.com:sub",
								"www.amazon.com:app_id",
								"www.amazon.com:user_id",
								"graph.facebook.com:app_id",
								"graph.facebook.com:id",
								"accounts.google.com:aud",
								"accounts.google.com:oaud",
								"accounts.google.com:sub",
								"aws:TagKeys",
								"aws:PrincipalTag/${TagKey}",
								"aws:RequestTag/${TagKey}",
								"sts:TransitiveTagKeys",
								"sts:SourceIdentity",
								"sts:RoleSessionName",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Decodes additional information about the authorization status of a request from an encoded message returned in response to an AWS request",
					Privilege:   "DecodeAuthorizationMessage",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Returns details about the access key id passed as a parameter to the request.",
					Privilege:   "GetAccessKeyInfo",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Returns details about the IAM identity whose credentials are used to call the API",
					Privilege:   "GetCallerIdentity",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Returns a set of temporary security credentials (consisting of an access key ID, a secret access key, and a security token) for a federated user",
					Privilege:   "GetFederationToken",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:PrincipalTag/${TagKey}",
								"aws:RequestTag/${TagKey}",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Returns a STS bearer token for an AWS root user, IAM role, or an IAM user",
					Privilege:   "GetServiceBearerToken",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Read",
					Description: "Returns a set of temporary security credentials (consisting of an access key ID, a secret access key, and a security token) for an AWS account or IAM user",
					Privilege:   "GetSessionToken",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Write",
					Description: "Grants permission to set a source identity on a STS session",
					Privilege:   "SetSourceIdentity",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user",
						},
						{
							ConditionKeys: []string{
								"sts:SourceIdentity",
								"aws:SourceIdentity",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
				{
					AccessLevel: "Tagging",
					Description: "Grants permission to add tags to a STS session",
					Privilege:   "TagSession",
					ResourceTypes: []ParliamentResourceType{
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "role",
						},
						{
							ConditionKeys:    []string{},
							DependentActions: []string{},
							ResourceType:     "user",
						},
						{
							ConditionKeys: []string{
								"aws:TagKeys",
								"aws:PrincipalTag/${TagKey}",
								"aws:RequestTag/${TagKey}",
								"sts:TransitiveTagKeys",
							},
							DependentActions: []string{},
							ResourceType:     "",
						},
					},
				},
			},
			Resources: []ParliamentResource{
				{
					Arn: "arn:${Partition}:iam::${Account}:role/${RoleNameWithPath}",
					ConditionKeys: []string{
						"aws:ResourceTag/${TagKey}",
					},
					Resource: "role",
				},
				{
					Arn:           "arn:${Partition}:iam::${Account}:user/${UserNameWithPath}",
					ConditionKeys: []string{},
					Resource:      "user",
				},
			},
			ServiceName: "AWS Security Token Service",
		},
	}

	return permissions
}
```