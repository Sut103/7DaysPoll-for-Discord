![linkedin_banner_image_1](https://github.com/Sut103/7DaysPoll-for-Discord/assets/18696845/df4b8411-1915-4d1b-81a2-381c2d8e5324)
# 7DaysPoll-for-Discord
Polling on 7 potential dates starting from the specified date.

## Manage Slash Commands
### Register
```
go run manage/* register
```

### Delete
```
go run manage/* delete
```

## Usage on AWS Lambda
### Run locally
```
docker build --platform linux/amd64 -t 7dayspoll:latest .
docker run -p 9000:8080 \
 --entrypoint /usr/local/bin/aws-lambda-rie \
 7dayspoll:latest /main
```

### Run on AWS
Please build the container image and push to your ECR.
Next creating the Lambda function with container image on your AWS account.
(Require environment variable: DISCORD_PUBLIC_KEY)