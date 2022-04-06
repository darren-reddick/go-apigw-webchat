# go-apigw-webchat

A simple example webchat API using AWS Apigateway 

<p align="left">
<img src="https://github.com/darren-reddick/go-apigw-webchat/actions/workflows/cicd.yml/badge.svg?branch=main">
</p>


[Usage](https://user-images.githubusercontent.com/57802771/160354279-fb4a817b-c990-49ea-98ad-169ab0d136c0.mov)


## :city_sunrise: Overview

This project provides a simple webchat service over a websocket API. Users can connect and send messages direct to other users or use a broadcast function to send a message to all users.

It is built solely using AWS services: Apigateway, Lambda, DynamoDB, Eventbridge and Cloudwatch.

<img src="https://user-images.githubusercontent.com/57802771/160415219-18ccd032-4e14-4a16-8e9d-014847127605.png" width="700">

## :factory: Deployment

### :page_facing_up: Prerequisites

- [Go (>=1.17)](https://go.dev/doc/install)
- Nodejs + npm >= 16
- An AWS account and IAM user with permissions to deploy

### Installing packages

Install node packages into project directory including serverless and wscat.

```
npm install
```

### :wrench: Provisioning the websocket service in AWS

A make target has been configured to build the binaries and deploy to AWS using the serverless framework. This will deploy to a named stage using the **STAGE** parameter. (stage is the equivalent of a deployment environment)

Configuration for a stage should be set in a file named [stage].yaml in the **environment** directory. See the **example.yaml** file as an example. Not including an environment file will simply take the defaults for each setting.

**The AWS IAM credentials for deploying should be set up in the environment to do this.**

**WARNING**: This will deploy the following resources to AWS which will incur costs:
- APIGateway websocket API
- Dynamodb table
- Lambda functions

```
# example to deploy to stage 'this'
make deploy STAGE=this
```

The stack output of the serverless deployment will list the **ServiceEndpointWebsocket** - make a note of this as it will be used to connect to the websocket API.

AWS IAM authentication for the API is now supported - to enable this see the **Authentication** section below.



## :notebook_with_decorative_cover: Usage

THe following examples use the wscat cli to use the webchat service. An alternative would be to use the piesocket online websocket client: https://www.piesocket.com/websocket-tester

### Connecting

Using wscat (ctrl-c to disconnect)

See the **Authentication** section below for connecting to an IAM authenticated API.

```
npx wscat -c [ServiceEndpointWebsocket]
```

The following actions can be done from within the wscat websocket session once connected:

### :loudspeaker: Broadcast a message

```
{"action":"broadcast","message":"Anybody home?"}
```

### :scroll: List connections (users)

```
{"action":"list"}
```

### :love_letter: Send a message to a user

```
{"action":"chat","message":"Yo!"}
```

## Authentication

Currently IAM authentication is supported for the websocket API for IAM users with the appropriate permissions on the API Gateway.

This can be enabled by updating the **iam_auth** setting to **true** in the environment file for the stage. E.g. **environment/poc.yaml** for the **poc** stage.

**NOTE**: if amending an existing API to enable authentication only the configuration of the API gets updated. The stage will need to be redeployed in the console for the authentication to activated for the deployment.

The initial connection route request will need to be signed with aws4 to authenticate. A script has been provided for this.

### Connecting

Using the script provided to sign the request using **aws4**. This requires the AWS IAM user to be set up as a named profile and its name set in environment variable **AWS_PROFILE**.

```
./scripts/wsconnect_auth.sh [stage name]
```

## :zombie: Removing the Deployment

```
# example to remove stage 'this'
make remove STAGE=this
```

## Logging

The zap structured logging package is used for logging. Debug logging can be enabled by setting the Lambda environment variable **LOG_LEVEL** to **DEBUG**. This can be set using the stage yaml files in the **environment** directory.

## API Routes

### connect

Initiate connection to the websocket API. Routes the lambda **connectHandler** function which persists the connection id to the data store and pushes a welcome message to an event bus.

### disconnect

Disconnect from the websocket API. Routes to the lambda **disconnectHandler** function which removes the connection id from the data store.

## broadcast

Broadcast a message to all connections (users). Routes to the lambda **broadcastHandler** function which queries the data store for all connection ids and sends the message to the associated connection.

## list

List all current connections to the websocket API. Routes to the lambda **listHandler** which returns a list of current connections to the websocket API by querying the data store to the requesting connection.

## chat

Send a message to a connection. Routes to the lambda **chatHandler** which sends the message to the connection id specified.



## Other Lambda functions

These functions are part of the service but arent directly connected from apigateway routes. 

### postHandler

Triggered by chat events from the chat eventbridge bus. Creates a message from the event and sends to the connected user. Initially this has been added to handle welcome messages which must be done asynchronously (a message cant be sent to the connection until it is established).

### manageHandler

Triggered by an eventbridge schedule to clean up connection ids from the database where the connection is no longer there.

# :pig: Known Issues

1. Serverless sometimes doesnt remove the API Gateway log group. This can be an issue if the project is redeployed with the same stage name and will fail with a message indicating that the log group already exists. To fix simply remove the log group manually using the console or the aws cli. (see Makefile **remove** target)






