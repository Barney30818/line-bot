# Line bot with serverless application

This line bot is used for reminding you to do something by schedulely, you can enhance it whatever you want such as notification, receive and response message, or call other APIs...

## Setup

1. Address all `TODO` comments.

2. Install [Node.js](https://nodejs.org/) and [Go](https://go.dev/).

3. Install and upgrade all packages to ensure your application is initialized with the latest package versions.  Note, this only need be done once.

       go get -t -u ./...
       go mod tidy

       npm update
       npx ncu -u

4. Commit changes

       git add .
       git commit

That's it, you are good to start coding!

For more on building AWS Lambdas with Go, see [AWS docs](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)

## Golang Linter Configuration

Change the .golangci.yml file to match your project's needs.

Execute the linter with:

    make lint

If you want to run the linter with auto-fixing, run:

    make lint-fix

## Build

Run the following command to build the application:

    make clean build

## Test

Testing setup and execution is left to the developer.

## Deploy


### Deploy to AWS on local

Run the following command to deploy the application:

    aws configure
    sls deploy --stage {stageName} --verbose

note: When deploy successfully, copy the url of `line-events` API, and paste to the Line developers -> Webhook URL

## API Specification

This document outlines the usage and specifications of two APIs. These APIs are designed for handling events from the LINE platform and broadcasting messages to users subscribed to a specific official account.

### /line-events

This API acts as a webhook for the LINE Bot, receiving and processing various events from the LINE platform.

#### Request Method

`POST`

### /notifications

This API allows for broadcasting messages to all users who have subscribed to the specific official account through a POST request.

#### Request Method

`POST`

#### Request Body

The body of the request should be in JSON format, containing a single field `message` which is a string representing the message you wish to broadcast. For example:

```json
{
  "message": "Did you take your fish oil today?\nPlease enter\n1. Yes\n2. No\nIf you haven't taken it or haven't responded, I will remind you again in 10 minutes."
}
