package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	// "github.com/aws/jsii-runtime-go"
)

type SlackBackendStackProps struct {
	awscdk.StackProps
}

//const dynamoDBTableName = "test_table_1"
var dynamoDBTableName string
var dynamoDBPartitionKey string
var dynamoDBTableNameEnvVar string

//const dynamoDBPartitionKey = "email"
//const dynamoDBTableNameEnvVar = "DYNAMODB_TABLE_NAME"

//const appRunnerServiceName = "dynamodb-apprunner-go-app"
//const appRunnerServicePort = "8080"

func init() {
	log.Println(".....running init() function.....")

	/*dynamoDBTableName = os.Getenv("DYNAMODB_TABLE_NAME")
	if dynamoDBTableName == "" {
		log.Fatal("missing env var DYNAMODB_TABLE_NAME")
	}

	dynamoDBPartitionKey = os.Getenv("DYNAMODB_PARTITION_KEY_ATTR")
	if dynamoDBPartitionKey == "" {
		log.Fatal("missing env var DYNAMODB_PARTITION_KEY_ATTR")
	}

	dynamoDBTableNameEnvVar = os.Getenv("APPRUNNER_DYNAMODB_TABLE_NAME_ENV_VAR")
	if dynamoDBTableNameEnvVar == "" {
		log.Fatal("missing env var APPRUNNER_DYNAMODB_TABLE_NAME_ENV_VAR")
	}*/

	slackSecret = os.Getenv(slackSecretEnvVar)
	if slackSecret == "" {
		log.Fatalf("missing env var %s\n", slackSecretEnvVar)
	}

	giphyAPIKey = os.Getenv(giphyAPIKeyEnvVar)
	if giphyAPIKey == "" {
		log.Fatalf("missing env var %s\n", giphyAPIKeyEnvVar)
	}
}

func main() {
	app := awscdk.NewApp(nil)

	NewSlackBackendStack(app, "AwsomeSlackBackendStack", &SlackBackendStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

var (
	slackSecret string
	giphyAPIKey string
)

const (
	functionName        = "awsome-slack-backend"
	functionBinaryName  = "awsome"
	functionZipFilePath = "../function.zip"
	slackSecretEnvVar   = "SLACK_SIGNING_SECRET"
	giphyAPIKeyEnvVar   = "GIPHY_API_KEY"
)

func NewSlackBackendStack(scope constructs.Construct, id string, props *SlackBackendStackProps) awscdk.Stack {

	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	//dynamoDBTable := awsdynamodb.NewTable(stack, jsii.String("dynamodb-test-table"), &awsdynamodb.TableProps{PartitionKey: &awsdynamodb.Attribute{Name: jsii.String(dynamoDBPartitionKey), Type: awsdynamodb.AttributeType_STRING}, TableName: jsii.String(dynamoDBTableName), RemovalPolicy: awscdk.RemovalPolicy_DESTROY})

	// environment variable for Lambda function
	lambdaEnvVars := &map[string]*string{slackSecretEnvVar: jsii.String(slackSecret), giphyAPIKeyEnvVar: jsii.String(giphyAPIKey)}

	//TODO include build function in CDK code itself

	// lambda function packaged as zip file
	//function := awslambda.NewFunction(stack, jsii.String("lambda-function"), &awslambda.FunctionProps{Runtime: awslambda.Runtime_GO_1_X(), Handler: jsii.String(functionBinaryName), Code: awslambda.AssetCode_FromAsset(jsii.String(functionZipFilePath), nil), Environment: lambdaEnvVars})

	//function := awslambda.NewDockerImageFunction(stack, jsii.String("func-docker"), &awslambda.DockerImageFunctionProps{FunctionName: jsii.String(functionName), Environment: lambdaEnvVars, Code: awslambda.DockerImageCode_FromImageAsset(jsii.String("."), &awslambda.AssetImageCodeProps{File: jsii.String("../../..")})})

	function := awslambda.NewDockerImageFunction(stack, jsii.String("func-docker"), &awslambda.DockerImageFunctionProps{FunctionName: jsii.String(functionName), Environment: lambdaEnvVars, Code: awslambda.DockerImageCode_FromImageAsset(jsii.String("../function"), nil)})

	funcURL := awslambda.NewFunctionUrl(stack, jsii.String("function-url"), &awslambda.FunctionUrlProps{AuthType: awslambda.FunctionUrlAuthType_NONE, Function: function})

	// policy to allow DynamoDB PutItem calls
	//dynamoDBPutItemPolicy := awsiam.NewPolicy(stack, jsii.String("policy"), &awsiam.PolicyProps{PolicyName: jsii.String("LambdaDynamoDBPutItemPolicy"), Statements: &[]awsiam.PolicyStatement{awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{Effect: awsiam.Effect_ALLOW, Actions: jsii.Strings("dynamodb:PutItem"), Resources: jsii.Strings(*dynamoDBTable.TableArn())})}})

	// attach the policy to an IAM role which is created during Lambda creation
	//function.Role().AttachInlinePolicy(dynamoDBPutItemPolicy)

	awscdk.NewCfnOutput(stack, jsii.String("Function URL"), &awscdk.CfnOutputProps{Value: funcURL.Url()})

	return stack
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
