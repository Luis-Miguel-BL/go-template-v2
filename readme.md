
# Go Boilerplate Essential Checklist

- [ ] **Project Structure & Documentation**
	- [X] Define a clean project structure (`cmd/`, `internal/`, etc.)
	- [ ] Document the project structure and conventions in the README
	- [ ] Add API endpoint documentation (Swagger/OpenAPI or similar)
	- [ ] Improve bootstrap logs

- [X] **Domain Layer**
	- [X] Implement the Domain Layer with aggregates, entities, value objects, and domain services
	- [X] Implement VO validations
	
- [X] **Application Layer**
	- [X] Implement application services / use cases
	- [X] Emit application events for side effects (integrations)

- [X] **Event Bus**
	- [X] Integrate a domain event bus for decoupled event handling
	- [X] Provide example domain events and subscribers

- [X] **Dependency Injection**
	- [X] Set up dependency injection using Uber Fx
	- [X] Register modules and lifecycle hooks

- [X] **Telemetry & Monitoring**
	- [X] Integrate OpenTelemetry for tracing and metrics
	- [X] Add New Relic instrumentation for APM
	- [X] Create examples for custom metrics
	- [X] Create examples for custom events
    - [X] Create examples for custom span

- [X] **AWS Integrations**
	- [X] Implement AWS DynamoDB connection and repository example
	- [X] Integrate AWS SSM for configuration management
    - [X] Add AWS SQS integration for messaging

- [X] **Serverless & Messaging Runners**
    - [X] Implement and run an SQS consumer
    - [X] Implement and run an AWS Lambda function

- [X] **HTTP Client & External Integrations**
	- [X] Add a standardized HTTP client (timeouts, retries, circuit breaker)
	- [X] Define an HTTP client interface for external APIs
	- [X] Implement an example external API integration
	- [X] Add observability to external calls (tracing, metrics, logs)
	- [X] Monitor external API latency, errors, and retries (OpenTelemetry + New Relic)

- [X] **Automated Testing**
	- [X] Add unit tests for core logic
	- [X] Add integration tests for API
	- [X] Add integration tests for Consumer
	- [X] Add integration tests for Lambda
	- [X] Set up test coverage reporting

- [ ] **Load Testing**
	- [ ] Provide load testing scripts or configuration (e.g., k6, Vegeta)

- [ ] **Profiling**
	- [ ] Provide profiling scripts or configuration (e.g., pprof)

- [ ] **CI/CD**
	- [ ] Set up CI pipeline for linting, testing, and building

- [X] **Utilities**
	- [X] Add utility functions (UUID generation, error handling, etc.)