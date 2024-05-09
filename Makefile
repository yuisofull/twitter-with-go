create-table:
	aws.exe dynamodb create-table --table-name tweet --attribute-definitions AttributeName=UserID,AttributeType=S AttributeName=TweetID,AttributeType=N --key-schema AttributeName=UserID,KeyType=HASH AttributeName=TweetID,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
delete-table:
	aws.exe dynamodb delete-table --table-name tweet
update-lambda:
	./list_tweet/build.sh
	./create_tweet/build.sh