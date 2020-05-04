S3 to Snowplow
==============

# Use case


# Concept
The idea is to leverage ops1 pipeline by pushing (Josh suggested go lambda) S3 Notifications (object key + byte size) as Snowplow events so they end up in ES.
At the same time we need to be able to get bands out of mt-configs2 repository (need to ask support to have the bands in clients repos, not only within the massive yaml) and extend sync-managed-service DAG to push the above Snowplow events into ops1. Github webhook might not fire off often enough and data might vanish.
Then with all the raw data in ES we can figure out a way of detecting the thresholds.
So the first step would be actually having a lambda posting notifications to ops1.

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

