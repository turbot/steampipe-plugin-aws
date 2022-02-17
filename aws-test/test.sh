#!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 BLACK="\e[30m"
 BOLDGREEN="\e[1;${GREEN}"
 ENDCOLOR="\e[0m"
 
 
# Define your function here
run_test () {
   echo -e "${BLACK}Running $1 ${ENDCOLOR}"
 if ! ./tint.js $1 >> output.txt
   then
    echo -e "${RED}Failed -> $1 ${ENDCOLOR}"
    echo $1 >> failed_tests.txt
  else
    echo -e "${BOLDGREEN}Passed -> $1 ${ENDCOLOR}"
    echo $1 >> passed_tests.txt
   fi
 }
 
 # output.txt - store output of each test
 # failed_tests.txt - names of failed test
 # passed_tests.txt names of passed test

 # removes files from previous test
# rm -rf output.txt failed_tests.txt passed_tests.txt
 date >> output.txt
 date >> failed_tests.txt
 date >> passed_tests.txt
run_test aws_accessanalyzer_analyzer
run_test aws_auditmanager_assessment
run_test aws_auditmanager_control
run_test aws_auditmanager_framework
run_test aws_cloudfront_distribution
run_test aws_cloudtrail_trail
run_test aws_codebuild_project
run_test aws_codepipeline_pipeline
run_test aws_config_configuration_recorder
run_test aws_config_conformance_pack
run_test aws_dax_cluster
run_test aws_dms_replication_instance
run_test aws_directory_service_directory
run_test aws_ec2_ami
run_test aws_ec2_application_load_balancer
run_test aws_ec2_classic_load_balancer
run_test aws_ec2_gateway_load_balancer
run_test aws_ec2_instance
run_test aws_ec2_load_balancer_listener
run_test aws_ec2_network_interface
run_test aws_ec2_network_load_balancer
run_test aws_ec2_reserved_instance
run_test aws_ec2_target_group
run_test aws_ec2_transit_gateway_vpc_attachment
run_test aws_ecs_cluster
run_test aws_efs_mount_target
run_test aws_eks_addon
run_test aws_eks_addon_version
run_test aws_eks_cluster
run_test aws_eks_identity_provider_config
run_test aws_elastic_beanstalk_application
run_test aws_elastic_beanstalk_environment
run_test aws_elasticache_cluster
run_test aws_elasticache_replication_group
run_test aws_elasticache_subnet_group
run_test aws_emr_cluster
run_test aws_emr_instance_group
run_test aws_fsx_file_system
run_test aws_glue_crawler
run_test aws_guardduty_finding
run_test aws_guardduty_ipset
run_test aws_guardduty_threat_intel_set
run_test aws_iam_access_advisor
run_test aws_iam_account_password_policy
run_test aws_iam_account_summary
run_test aws_identitystore_group
run_test aws_identitystore_user
run_test aws_kinesis_consumer
run_test aws_kinesis_firehose_delivery_stream
run_test aws_macie2_classification_job
run_test aws_media_store_container
run_test aws_organizations_account
run_test aws_rds_db_cluster
run_test aws_rds_db_cluster_snapshot
run_test aws_rds_db_event_subscription
run_test aws_rds_db_instance
run_test aws_rds_db_snapshot
run_test aws_rds_db_subnet_group
run_test aws_redshift_cluster
run_test aws_redshift_snapshot
run_test aws_redshift_subnet_group
run_test aws_route53_resolver_endpoint
run_test aws_s3_access_point
run_test aws_s3_bucket
run_test aws_securityhub_hub
run_test aws_securityhub_product
run_test aws_securityhub_standards_subscription
run_test aws_servicequotas_default_service_quota
run_test aws_servicequotas_service_quota
run_test aws_sfn_state_machine_execution
run_test aws_sfn_state_machine_execution_history
run_test aws_ssm_association
run_test aws_ssm_managed_instance
run_test aws_ssm_managed_instance_compliance
run_test aws_ssm_parameter
run_test aws_ssoadmin_instance
run_test aws_ssoadmin_managed_policy_attachment
run_test aws_ssoadmin_permission_set
run_test aws_vpc
run_test aws_vpc_egress_only_internet_gateway
run_test aws_vpc_endpoint
run_test aws_vpc_endpoint_service
run_test aws_vpc_flow_log
run_test aws_vpc_internet_gateway
run_test aws_vpc_nat_gateway
run_test aws_vpc_network_acl
run_test aws_vpc_route
run_test aws_vpc_route_table
run_test aws_vpc_security_group
run_test aws_vpc_subnet
run_test aws_vpc_vpn_connection
run_test aws_vpc_vpn_gateway
run_test aws_wafv2_rule_group
run_test aws_wafv2_web_acl

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt