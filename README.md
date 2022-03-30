# go-api-webchat

A simple example webchat API using AWS Apigateway 

<p align="left">
<img src="https://github.com/darren-reddick/go-apigw-webchat/actions/workflows/cicd.yml/badge.svg?branch=main">
</p>

## :city_sunrise: Overview

This project provides a simple webchat service over a websocket API. Users can connect and send messages direct to other users or use a broadcast function to send a message to all users.

It is built solely on AWS services: Apigateway, Lambda, DynamoDB, Eventbridge and Cloudwatch.


## :factory: Deployment

### :page_facing_up: Prerequisites

- [Go (>=1.17)](https://go.dev/doc/install)
- [Serverless Framework 2](https://www.serverless.com/framework/docs/getting-started)
- Nodejs + npm (for serverless framework)
- An AWS account and IAM user with permissions to deploy
- wscat (optional for testing)

### Installing packages for serverless

```
npm install
```

### :wrench: Provisioning the websocket service in AWS

A make target has been configured to build the binaries and deploy to AWS using the serverless framework. **The AWS IAM credentials for deploying should be set up in the environment to do this.**

**WARNING**: This will deploy the following resources to AWS which will incur costs:
- APIGateway websocket API
- Dynamodb table
- Lambda functions

```
make deploy
```

The stack output of the serverless deployment will list the **ServiceEndpointWebsocket** - make a note of this as it will be used to connect to the websocket API.


## :notebook_with_decorative_cover: Usage

### Connecting

Using wscat (ctrl-c to disconnect)

```
wscat -c [ServiceEndpointWebsocket]
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

### connectHandler

Handles the initiation of the websocket connection from the client via apigateway. This involves persisting the connection id to the data store and pushes a welcome message to an event bus.

### postHandler

Triggered by chat events from the chat eventbridge bus. Creates a message from the event and sends to the connected user. Initially this has been added to handle welcome messages which must be done asynchronously (a message cant be sent to the connection until it is established).

### manageHandler

Triggered by an eventbridge schedule to clean up connection ids from the database where the connection is no longer there.





