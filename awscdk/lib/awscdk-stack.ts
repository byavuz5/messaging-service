import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as iam from 'aws-cdk-lib/aws-iam'
// import * as sqs from 'aws-cdk-lib/aws-sqs';

export class AwscdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Set ecr image properties
    const ecrImageCodeProps: lambda.EcrImageCodeProps = {
      tagOrDigest: 'latest',
    };

   
    // Create api gateway
    const api = new apigateway.RestApi(this, 'messaging_api', {
      description: 'messaging api gateway',
      deployOptions: {
        stageName: 'prod',
      },
    });

    // Create dynamoDB tables
    
    // Create users tables with dynamoDB
    const usersTable = new dynamodb.Table(this, 'users', {
      tableName:"users", 
      partitionKey: { name: 'username', type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PROVISIONED, 
      readCapacity: 5,
      writeCapacity: 5,
      tableClass: dynamodb.TableClass.STANDARD,
    });

    // Create messages tables with dynamoDB
    const messagesTable = new dynamodb.Table(this, 'messages', {
      tableName:"messages",  
      partitionKey: { name: 'sender', type: dynamodb.AttributeType.STRING },
      sortKey: {name: 'createdAt', type: dynamodb.AttributeType.STRING},
      billingMode: dynamodb.BillingMode.PROVISIONED, 
      readCapacity: 5,
      writeCapacity: 5,
      tableClass: dynamodb.TableClass.STANDARD,
    });


     // Create system_logs tables with dynamoDB
     const systemLogsTable = new dynamodb.Table(this, 'system_logs', { 
      tableName:"system_logs", 
      partitionKey: { name: 'service_name', type: dynamodb.AttributeType.STRING },
      sortKey: {name: 'createdAt', type: dynamodb.AttributeType.STRING},
      billingMode: dynamodb.BillingMode.PROVISIONED, 
      readCapacity: 5,
      writeCapacity: 5,
      tableClass: dynamodb.TableClass.STANDARD,
    });

     // Create activity_logs tables with dynamoDB
     const activityLogsTable = new dynamodb.Table(this, 'activity_logs', {
      tableName:"activity_logs",  
      partitionKey: { name: 'username', type: dynamodb.AttributeType.STRING },
      sortKey: {name: 'createdAt', type: dynamodb.AttributeType.STRING},
      billingMode: dynamodb.BillingMode.PROVISIONED, 
      readCapacity: 5,
      writeCapacity: 5,
      tableClass: dynamodb.TableClass.STANDARD,
    });


    // Create lambda functions

    // Create ecr repository for login service docker image
    var loginRepo = cdk.aws_ecr.Repository.fromRepositoryName(this, "loginRepository","login")
    
     // Create login lambda service with docker image
    const loginService = new lambda.DockerImageFunction(this, 'login-service', {
      functionName:"login-service",
      code: lambda.DockerImageCode.fromEcr(loginRepo,ecrImageCodeProps),
    });

    //Add login lambda function policy for Query access users table  

    loginService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:Query'],
        resources: [usersTable.tableArn]
      })
    );
    
    //Add login lambda function policy for PutItem access systemLogsTable table  
    loginService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [systemLogsTable.tableArn]
      })
    );

    //Add login lambda function policy for PutItem access activityLogsTable table  
    loginService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [activityLogsTable.tableArn]
      })
    );


    // Create ecr repository for create-account service docker image
    var createAccountRepo = cdk.aws_ecr.Repository.fromRepositoryName(this, "createAccountRepository","create-account")

    // Create create-account lambda service with docker image
    const createAccountService = new lambda.DockerImageFunction(this, 'create-account-service', {
      functionName:"create-account-service",
      code: lambda.DockerImageCode.fromEcr(createAccountRepo,ecrImageCodeProps),
    });

    //Add createAccount lambda function policy for PutItem access users table  
    createAccountService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [usersTable.tableArn]
      })
    );
    
    //Add createAccount lambda function policy for PutItem access systemLogsTable table  
    createAccountService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [systemLogsTable.tableArn]
      })
    );
    
    

    // Create ecr repository for getAllMessage service docker image
    var getAllMessagesRepo = cdk.aws_ecr.Repository.fromRepositoryName(this, "getAllMessagesRepository","get-all-messages")
    
     // Create getAllMessage lambda service with docker image
    const getAllMessagesService = new lambda.DockerImageFunction(this, 'get-all-messages', {
      functionName:"get-all-messages",
      code: lambda.DockerImageCode.fromEcr(getAllMessagesRepo,ecrImageCodeProps),
    });

    //Add getAllMessages lambda function policy for Scan access users table  
    getAllMessagesService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:Scan'],
        resources: [messagesTable.tableArn]
      })
    );

    //Add getAllMessages lambda function policy for PutItem access systemLogsTable table  
    getAllMessagesService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [systemLogsTable.tableArn]
      })
    );

    //Add getAllMessages lambda function policy for PutItem access activityLogsTable table  
    getAllMessagesService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [activityLogsTable.tableArn]
      })
    );


      
    // Create ecr repository for getContactMessage service docker image
    var getContactMessagesRepo = cdk.aws_ecr.Repository.fromRepositoryName(this, "getContactMessagesRepository","get-contact-messages")
    
     // Create getContactMessage lambda service with docker image
    const getContactMessagesService = new lambda.DockerImageFunction(this, 'get-contact-messages', {
      functionName:"get-contact-messages",
      code: lambda.DockerImageCode.fromEcr(getContactMessagesRepo,ecrImageCodeProps),
    });

    //Add getContactMessages lambda function policy for Scan access users table  
    getContactMessagesService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:Scan'],
        resources: [messagesTable.tableArn]
      })
    );

    //Add getContactMessages lambda function policy for PutItem access systemLogsTable table  
    getContactMessagesService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [systemLogsTable.tableArn]
      })
    );

    //Add getContactMessages lambda function policy for PutItem access activityLogsTable table  
    getContactMessagesService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [activityLogsTable.tableArn]
      })
    );
      
    // Create ecr repository for sendMessage service docker image
    var sendMessageRepo = cdk.aws_ecr.Repository.fromRepositoryName(this, "sendMessageRepository","send-message")
    
     // Create sendMessage lambda service with docker image
    const sendMessageService = new lambda.DockerImageFunction(this, 'send-message', {
      functionName:"send-message",
      code: lambda.DockerImageCode.fromEcr(sendMessageRepo,ecrImageCodeProps),
    });

    //Add sendMessage lambda function policy for PutItem access users table  
    sendMessageService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [messagesTable.tableArn]
      })
    );

    //Add sendMessage lambda function policy for PutItem access systemLogsTable table  
    sendMessageService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [systemLogsTable.tableArn]
      })
    );

    //Add sendMessage lambda function policy for PutItem access activityLogsTable table  
    sendMessageService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [activityLogsTable.tableArn]
      })
    );

   

    // Create ecr repository for getSystemLogs service docker image
    var getSystemLogsRepo = cdk.aws_ecr.Repository.fromRepositoryName(this, "getSystemLogsRepository","get-system-logs")
    
     // Create getSystemLogs lambda service with docker image
    const getSystemLogsService = new lambda.DockerImageFunction(this, 'get-system-logs', {
      functionName:"get-system-logs",
      code: lambda.DockerImageCode.fromEcr(getSystemLogsRepo,ecrImageCodeProps),
     
    });

    //Add getSystemLogs lambda function policy for Query access users table
    getSystemLogsService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:Query'],
        resources: [systemLogsTable.tableArn]
      })
    );

    //Add getSystemLogs lambda function policy for PutItem access systemLogsTable table  
    getSystemLogsService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [systemLogsTable.tableArn]
      })
    );

   

    

    // Create ecr repository for getActivityLogs service docker image
    var getActivityLogsRepo = cdk.aws_ecr.Repository.fromRepositoryName(this, "getActivityLogsRepository","get-activity-logs")
    
     // Create getActivityLogs lambda service with docker image
    const getActivityLogsService = new lambda.DockerImageFunction(this, 'get-activity-logs', {
      functionName:"get-activity-logs",
      code: lambda.DockerImageCode.fromEcr(getActivityLogsRepo,ecrImageCodeProps),
      
    });
    
    //Add getActivityLogs lambda function policy for Query access users table
    getActivityLogsService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:Query'],
        resources: [activityLogsTable.tableArn]
      })
    );

    //Add getActivityLogs lambda function policy for PutItem access systemLogsTable table  
    getActivityLogsService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [systemLogsTable.tableArn]
      })
    );

    //Add getActivityLogs lambda function policy for PutItem access activityLogsTable table  
    getActivityLogsService.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        resources: [activityLogsTable.tableArn]
      })
    );
    
    //Create api resoruces

    //Add login resource
    const loginResource = api.root.addResource('login');
    
    

    //login resource add login service
    loginResource.addMethod(
      'POST',
      new apigateway.LambdaIntegration(loginService, {proxy: false,integrationResponses: [{
        statusCode: '200',
      }],}
        
        
        ),
      {
        methodResponses: [{
          statusCode: "200",
          responseModels: { "application/json": apigateway.Model.EMPTY_MODEL }
        }],
      }
    );

    
    
    //Add createAccount resource
    const createAccountResource = api.root.addResource('createAccount');
    
    //createAccount resource add createAccount service
    createAccountResource.addMethod(
      'POST',
      new apigateway.LambdaIntegration(createAccountService, {proxy: false,integrationResponses: [{
        statusCode: '200',
      }],}
        
        
        ),
      {
        methodResponses: [{
          statusCode: "200",
          responseModels: { "application/json": apigateway.Model.EMPTY_MODEL }
        }],
      }
    );
      

      
    //Add getAllMessages resource
    const getAllMessagesResource = api.root.addResource('getAllMessages');
    
    //getAllMessages resource add getAllMessages service
    getAllMessagesResource.addMethod(
      'POST',
      new apigateway.LambdaIntegration(getAllMessagesService, {proxy: false,integrationResponses: [{
        statusCode: '200',
      }],}
        
        
        ),
      {
        methodResponses: [{
          statusCode: "200",
          responseModels: { "application/json": apigateway.Model.EMPTY_MODEL }
        }],
      }
    );

      

    //Add getContactMessages resource
    const getContactMessagesResource = api.root.addResource('getContactMessages');
    
    //getContactMessages resource add getContactMessages service
    getContactMessagesResource.addMethod(
      'POST',
      new apigateway.LambdaIntegration(getContactMessagesService, {proxy: false,integrationResponses: [{
        statusCode: '200',
      }],}
        
        
        ),
      {
        methodResponses: [{
          statusCode: "200",
          responseModels: { "application/json": apigateway.Model.EMPTY_MODEL }
        }],
      }
    );


    //Add sendMessage resource
    const sendMessageResource = api.root.addResource('sendMessage');
    
    //sendMessage resource add sendMessage service
    sendMessageResource.addMethod(
      'POST',
      new apigateway.LambdaIntegration(sendMessageService, {proxy: false,integrationResponses: [{
        statusCode: '200',
      }],}
        
        
        ),
      {
        methodResponses: [{
          statusCode: "200",
          responseModels: { "application/json": apigateway.Model.EMPTY_MODEL }
        }],
      }
    );


    //Add getSystemLogs resource
    const getSystemLogsResource = api.root.addResource('getSystemLogs');
    
    //getSystemLogs resource add sendSystemLog service
    getSystemLogsResource.addMethod(
      'POST',
      new apigateway.LambdaIntegration(getSystemLogsService, {proxy: false,integrationResponses: [{
        statusCode: '200',
      }],}
        
        
        ),
      {
        methodResponses: [{
          statusCode: "200",
          responseModels: { "application/json": apigateway.Model.EMPTY_MODEL }
        }],
      }
    );



    //Add getActivityLogs resource
    const getActivityLogsResource = api.root.addResource('getActivityLogs');
    
    //getActivityLogs resource add sendActivityLog service
    getActivityLogsResource.addMethod(
      'POST',
      new apigateway.LambdaIntegration(getActivityLogsService, {proxy: false,integrationResponses: [{
        statusCode: '200',
      }],}
        
        
        ),
      {
        methodResponses: [{
          statusCode: "200",
          responseModels: { "application/json": apigateway.Model.EMPTY_MODEL }
        }],
      }
    );

    // Create an Output for the API URL
    new cdk.CfnOutput(this, 'apiUrl', {value: api.url});
  }
}
