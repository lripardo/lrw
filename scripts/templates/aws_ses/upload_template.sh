#!/bin/sh
if [ -z "$1" ]
  then
    echo "No method was found. Ex: ./upload_template.sh create-template"
    exit 1
fi

if [ -z "$2" ]
  then
    echo "No file name for template was found. Ex: ./upload_template.sh update-template verify_template.json"
    exit 1
fi

# need to put aws credentials on /home/user/.aws/ folder. More info here: https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2-docker.html
docker run --rm -it -v ~/.aws:/root/.aws -v $(pwd):/files amazon/aws-cli ses $1 --cli-input-json fileb:///files/$2