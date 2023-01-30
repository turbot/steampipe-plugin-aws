 #!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 BLACK="\e[30m"
 BOLDGREEN="\e[1;${GREEN}"
 ENDCOLOR="\e[0m"

# Define your function here
run_test () {
   echo -e "${GREEN}Running $1 ${ENDCOLOR}"
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



run_test aws_sagemaker_endpoint_configuration
run_test aws_sagemaker_model
run_test aws_ssm_association
run_test aws_ssm_managed_instance
run_test aws_ssm_managed_instance_compliance
run_test aws_vpc
run_test aws_vpc_egress_only_internet_gateway
run_test aws_vpc_endpoint
run_test aws_vpc_endpoint_service
run_test aws_vpc_flow_log
run_test aws_vpc_internet_gateway
run_test aws_vpc_nat_gateway
run_test aws_vpc_network_acl
run_test aws_vpc_peering_connection
run_test aws_vpc_route
run_test aws_vpc_route_table
run_test aws_vpc_security_group
run_test aws_vpc_subnet
run_test aws_vpc_vpn_connection
run_test aws_vpc_vpn_gateway

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt