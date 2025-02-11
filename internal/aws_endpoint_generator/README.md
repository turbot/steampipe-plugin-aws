
# **AWS Endpoint Generator**

## **Overview**
The **AWS Endpoint Generator** is a Go-based tool that automatically fetches AWS service endpoint data from the **AWS SDK for Go v2** and generates Go files containing:

1. **AWS-supported service endpoints** across partitions and regions.
2. **Service ID constants** to facilitate programmatic access to service metadata.

This tool helps in determining **AWS service availability** across different regions and partitions.

---

## **How It Works**
The tool retrieves endpoint data from the following AWS SDK JSON file:  
🔗 **[AWS Endpoints JSON](https://raw.githubusercontent.com/aws/aws-sdk-go-v2/master/codegen/smithy-aws-go-codegen/src/main/resources/software/amazon/smithy/aws/go/codegen/endpoints.json)**  

It then **parses the JSON**, extracts relevant data, and generates two Go files:
- **`aws/endpoint_list_gen.go`** → Contains partition, region, and service endpoint definitions.
- **`aws/endpoint_service_ids_gen.go`** → Contains service ID constants for all AWS services.

---

## **Installation**
### **1. Move to the dir steampipe-plugin-aws**
```sh
cd steampipe-plugin-aws
```

### **2. Install Dependencies**
Ensure Go modules are up to date:
```sh
go mod tidy
```

### **3. Run the Endpoint Generator script**

```sh
go run internal/aws_endpoint_generator/aws_supported_endpoints.go 
```

---

## **Generated Files**
The tool generates the following files inside `aws/`:

### **1. `endpoint_list_gen.go`** (AWS Service Endpoints)
Contains structured **Go definitions** of AWS partitions, regions, and service endpoints.

Example:
```go
var AWSPartition = Partition{
	ID:          "aws",
	Name:        "AWS Standard",
	DNSSuffix:   "amazonaws.com",
	Regions: map[string]Region{
		"us-east-1": {ID: "us-east-1", Description: "US East (N. Virginia)"},
		"us-west-2": {ID: "us-west-2", Description: "US West (Oregon)"},
	},
}
```

### **2. `endpoint_service_ids_gen.go`** (Service ID Constants)
Defines **constants** for AWS services.

Example:
```go
const AWS_S3_SERVICE_ID = "s3"
const AWS_EC2_SERVICE_ID = "ec2"
const AWS_LAMBDA_SERVICE_ID = "lambda"
```