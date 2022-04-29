package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type SlackBackendStackProps struct {
	awscdk.StackProps
}

func init() {

	slackSecret = os.Getenv(slackSecretEnvVar)
	if slackSecret == "" {
		log.Fatalf("missing environment variable %s\n", slackSecretEnvVar)
	}

	giphyAPIKey = os.Getenv(giphyAPIKeyEnvVar)
	if giphyAPIKey == "" {
		log.Fatalf("missing environment variable %s\n", giphyAPIKeyEnvVar)
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
	functionName       = "awsome-slack-backend"
	functionBinaryName = "awsome"
	//functionZipFilePath = "../function.zip"
	slackSecretEnvVar = "SLACK_SIGNING_SECRET"
	giphyAPIKeyEnvVar = "GIPHY_API_KEY"
)

func NewSlackBackendStack(scope constructs.Construct, id string, props *SlackBackendStackProps) awscdk.Stack {

	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// environment variable for Lambda function
	lambdaEnvVars := &map[string]*string{slackSecretEnvVar: jsii.String(slackSecret), giphyAPIKeyEnvVar: jsii.String(giphyAPIKey)}

	function := awslambda.NewDockerImageFunction(stack, jsii.String("awsome-func-docker"), &awslambda.DockerImageFunctionProps{FunctionName: jsii.String(functionName), Environment: lambdaEnvVars, Code: awslambda.DockerImageCode_FromImageAsset(jsii.String("../function"), nil)})

	funcURL := awslambda.NewFunctionUrl(stack, jsii.String("awsome-func-url"), &awslambda.FunctionUrlProps{AuthType: awslambda.FunctionUrlAuthType_NONE, Function: function})

	awscdk.NewCfnOutput(stack, jsii.String("Function URL"), &awscdk.CfnOutputProps{Value: funcURL.Url()})

	return stack
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
