.PHONY: run-api run-worker test create-dynamo-table create-sqs-queue coverage-html coverage-cli


APP_AWS_ENDPOINT=http://localhost:4566
APP_AWS_REGION=sa-east-1
APP_AWS_SSM_LOAD_FROM_SSM=false
APP_ENVIRONMENT=local

AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
AWS_SESSION_TOKEN=test

LOCAL_ENV = \
	APP_AWS_ENDPOINT=$(APP_AWS_ENDPOINT) \
	APP_AWS_REGION=$(APP_AWS_REGION) \
	APP_AWS_SSM_LOAD_FROM_SSM=$(APP_AWS_SSM_LOAD_FROM_SSM) \
	APP_ENVIRONMENT=$(APP_ENVIRONMENT) \
	AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
	AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
	AWS_SESSION_TOKEN=$(AWS_SESSION_TOKEN)

run-api:
	$(LOCAL_ENV) go run cmd/api/main.go

run-worker:
	$(LOCAL_ENV) go run cmd/worker/main.go

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
	go test ./... -coverprofile=test/coverage/coverage.out -coverpkg=./...   

coverage-html: test
	go tool cover -html=test/coverage/coverage.out -o test/coverage/coverage.html
	open test/coverage/coverage.html

coverage-cli: test
	go tool cover -func=test/coverage/coverage.out
	

load-test:
	go run test/load/create_lead.go