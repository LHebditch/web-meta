package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/jsii-runtime-go"

	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	bundlingOptions := &awslambdago.BundlingOptions{
		GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w"`)},

	}
	notFound := awslambdago.NewGoFunction(stack, jsii.String("WebMetaApi_NotFound"), &awslambdago.GoFunctionProps{
		FunctionName: jsii.String("web-meta-api-not-found"),
		Entry:    jsii.String("../handlers/notfound"),
		Bundling: bundlingOptions,
	})

	webscrapper := awslambdago.NewGoFunction(stack, jsii.String("WebMetaApi_ScrapeMeta"), &awslambdago.GoFunctionProps{
		FunctionName: jsii.String("web-meta-api-scrape-meta"),
		Entry:    jsii.String("../handlers/webmeta/get/"),
		Bundling: bundlingOptions,
	})

	webMetaApi := awsapigateway.NewLambdaRestApi(stack, jsii.String("web-meta-api"), &awsapigateway.LambdaRestApiProps{
		Proxy: jsii.Bool(false),
		CloudWatchRole: jsii.Bool(false),
		Handler: notFound,
	})
	// apiResourceOpts := new(awsapigateway.ResourceOptions)
	apiLambdaOpts := new(awsapigateway.LambdaIntegrationOptions)
	apiMethodOpts := new(awsapigateway.MethodOptions)

	webMetaApi.Root().AddMethod(jsii.String("POST"), awsapigateway.NewLambdaIntegration(webscrapper, apiLambdaOpts), apiMethodOpts)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "WebMeta", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
