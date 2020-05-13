S3 to Snowplow
==============

# Use case
Send notifications whenever data input is larger than estimated processor throughput.

# Concept
Snowplow ops1 pipeline can consume events for:
- s3 incoming data (via lambda s3 notifications)
- data processor size (via mt-configs2 repository or via bands.yml file, gh webhook discouraged due to possible data vanishing between commits)
The data is then put into ES where it is cross-referenced.
The idea is to leverage ops1 pipeline by pushing (Josh suggested go lambda) S3 Notifications (object key + byte size) as Snowplow events so they end up in ES.

# Running
Required components:
- Snowplow (ie. micro)
- AWS S3, Lambda (ie. Localstack w/ `SERVICES=serverless,lambda,s3`)

0. Fire up with docker: 
```
TMPDIR=/private$TMPDIR docker-compose up
```
1. Deploy to lambda & s3 bucket
```
aws --endpoint-url http://localhost:4574 lambda create-function \
  --function-name s3-events --runtime go1.x \
  --zip-file fileb://function.zip --handler main \
  --role arn:aws:iam::123456:role/irrelevant \
  --environment Variables="{EVENT_SCHEMA='TBD', COLLECTOR_URI='TBD'}"

aws --endpoint-url=http://localhost:4572 s3 mb my-bucket
```
2. Setup notifications
```
aws --endpoint-url=http://localhost:4572 s3api put-bucket-notification-configuration --bucket my-bucket --notification-configuration file://config/notification.json
```
2. Put anything to s3
```
touch any.txt
aws --endpoint-url=http://localhost:4572 s3 cp any.txt  s3://my-bucket/any.txt
```
3. Get an event in micro

# Links
- Jobs https://docs.google.com/spreadsheets/d/1AlKtU87pzSV3fnBsncEgquS6_05BKbny6uh-Vlg5CNQ/edit#gid=1733119665
- Processor bands https://docs.google.com/spreadsheets/d/1RD1RYUoWRRTBBNTzA7pQNyA5jro7Q1obe5w1E9oM3sg/edit#gid=769662015
- Processor bands sizes https://github.com/snowplow-devops/support-tools/blob/master/support/basebands.yml
- Bands submitting can happen like here https://github.com/snowplow-proservices/snowplow-techops-pipeline/blob/master/schedules/sync-managed-service.json https://github.com/snowplow-proservices/mt-scripts/blob/master/snowplowtechops/dags/sync-managed-service.factfile (this is where it is implemented: https://github.com/snowplow-devops/sync-managed-service/blob/master/sync_managed_service/sync.py#L362-L365)
- On go-lambda side it might be a good idea to implement flush sending, but might not be necessary https://github.com/snowplow-devops/metrics-relay/blob/master/tracking.go#L156-L188
- Go lambda handler https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
- Go tracker: https://github.com/snowplow/snowplow-golang-tracker
