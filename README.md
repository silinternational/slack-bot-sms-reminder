# Slack Bot SMS Reminders
This project is primarily an experiment and learning tool for playing with slack bots and AWS Lambda. Please do not
actually use this, it's just not that good :-)

## Setup
1. Create a new app at [https://api.slack.com](https://api.slack.com)
2. If you do not already have a Nexmo account, register for one at [https://www.nexmo.com](https://www.nexmo.com)
3. Copy `aws.env.example` to `aws.env`
4. Edit `aws.env` and populate with appropriate values
5. Make sure you have Docker and Docker Compose installed
6. Run `make deploy` - This will compile the Go binaries and run `serverless deploy` to create an API Gateway and deploy
the functions to Lambda
7. From the output, copy the endpoint, something like `https://1j876uqzle.execute-api.us-east-1.amazonaws.com/dev/sms` 
8. In your Slack App configuration, add a Slash command for `/sms` and paste in the endpoint url you just copied. 
9. Install the app into a Slack workspace you have permission to do so in
10. Test the integration by typing in `/sms [yourphonenumber] 1m your message` and you should get a response telling 
you when it was scheduled for and then in that much time you should also receive the text message. Note that if you 
are using a free trial Nexmo account you wont be able to text any numbers than the one you registered with. 

## About
This app was created as a demo for a presentation given to the Charlotte Golang Meetup group in August 2018. The slides 
for that talk are available at [https://bit.ly/golang-clt-serverless](https://bit.ly/golang-clt-serverless)