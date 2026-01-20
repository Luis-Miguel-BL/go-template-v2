.PHONY: run-api run-worker test create-dynamo-table create-sqs-queue

run-api:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

create-dynamo-table:
	aws dynamodb create-table \
		--table-name lead_single_table \
		--attribute-definitions \
			AttributeName=PK,AttributeType=S \
			AttributeName=SK,AttributeType=S \
		--key-schema \
			AttributeName=PK,KeyType=HASH \
			AttributeName=SK,KeyType=RANGE \
		--billing-mode PAY_PER_REQUEST \
		--endpoint-url http://localhost:4566 \
		--region sa-east-1

create-sqs-queue:
	aws sqs create-queue \
		--queue-name lead-created \
		--endpoint-url http://localhost:4566 \
		--region sa-east-1

get-sqs-queue-url:
	aws sqs get-queue-url \
	  --queue-name lead-created \
	  --endpoint-url http://localhost:4566 \
	  --region sa-east-1

test:
	go test -failfast ./... -v

