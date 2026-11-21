default: test

# A failing go test must fail make even when its output goes through tee (DEBUG=1)
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -c

# With DEBUG=1 the full output goes to the log file and the terminal shows only the test results, the skip and failure messages
LOGFILTER = grep --line-buffered -E '^(=== RUN|--- |PASS|FAIL|ok |[?] |panic|exit status)|^[[:space:]]+[A-Za-z0-9_]+[.]go:[0-9]+:' || true

# Load environment variables from the .env file if it exists (copy .env.sample to .env)
# Write the values without quotes; a literal $ in a value must be written as $$.
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# A NAME or DEBUG in the shell environment (WSL sets NAME to the host name) is not a request: only the command line counts
ifeq ($(origin NAME),environment)
    NAME :=
endif
ifeq ($(origin DEBUG),environment)
    DEBUG :=
endif

# Run the acceptance tests (the TestAcc* tests) against the Crosswork cluster of the .env file
# Usage: make test [NAME="Definition Name"] [DEBUG=1]
# The offline unit tests are not part of it: use make unit
# DEBUG=1 sets TF_LOG=DEBUG, clears any TF_LOG_PATH and writes the full output, Terraform log included, to test-output.log; the terminal shows the test results only
# NAME must be a definition name (e.g., "Users", "Network Devices") matching the name: field in gen/definitions/*.yaml
# Tests that need a tag skip themselves when it is not set in .env
.PHONY: test
test:
	@echo "========================================="
	@echo "Running acceptance tests..."
	@echo "========================================="
	@if [ -z "$(CW_URL)" ] || [ -z "$(CW_USERNAME)" ] || [ -z "$(CW_PASSWORD)" ]; then \
		echo "SKIPPED: CW_URL, CW_USERNAME and CW_PASSWORD are not all configured"; \
		echo "To enable tests, copy .env.sample to .env and fill in CW_URL, CW_USERNAME and CW_PASSWORD"; \
	else \
		TEST_NAME="TestAcc"; \
		if [ -n "$(NAME)" ]; then \
			MATCH=$$(grep -l '^name: $(NAME)$$' gen/definitions/*.yaml 2>/dev/null | head -1); \
			if [ -n "$${MATCH}" ]; then \
				CAMEL=$$(echo "$(NAME)" | awk '{for(i=1;i<=NF;i++) $$i=toupper(substr($$i,1,1)) substr($$i,2)}1' | tr -d ' '); \
				case "$${CAMEL}" in GoldenConfig*) CAMEL=GoldenConfig;; esac; \
				TEST_NAME="TestAcc.*Cw$${CAMEL}"; \
				echo "Resolved definition '$(NAME)' to test pattern: $${TEST_NAME}"; \
			else \
				echo "ERROR: No definition found matching '$(NAME)' in gen/definitions/*.yaml"; \
				exit 1; \
			fi; \
		fi; \
		$(if $(DEBUG),echo "Debug mode enabled - the full output is written to test-output.log and the terminal shows the test results only";) \
		if [ -n "$${TEST_NAME}" ]; then echo "Running tests matching: $${TEST_NAME}"; fi; \
		TF_ACC=1 \
		$(if $(DEBUG),env -u TF_LOG_PATH TF_LOG=DEBUG) \
		go test -v $${TEST_NAME:+-run "$${TEST_NAME}"} $(TESTARGS) -count 1 -timeout 120m ./internal/provider $(if $(DEBUG),2>&1 | tee test-output.log | { $(LOGFILTER); }); \
	fi

# Same as test, so that "make testacc" keeps working
.PHONY: testacc
testacc: test

# Run the offline unit tests, which need no cluster; the acceptance tests skip themselves because TF_ACC is not set
# Usage: make unit [TESTARGS="-run TestPool"] [DEBUG=1]
# DEBUG=1 runs the tests verbosely with TF_LOG=DEBUG and writes the full output to unit-output.log; the terminal shows the test results only
# It clears any TF_LOG_PATH so that the Terraform log goes into the same file
.PHONY: unit
unit:
	$(if $(DEBUG),@echo "Debug mode enabled - the full output is written to unit-output.log and the terminal shows the test results only")
	@env -u TF_ACC $(if $(DEBUG),-u TF_LOG_PATH TF_LOG=DEBUG) go test $(if $(DEBUG),-v) ./... -count 1 $(TESTARGS) $(if $(DEBUG),2>&1 | tee unit-output.log | { $(LOGFILTER); })

# Regenerate the code, the examples and the documentation from the definitions
# Usage: make gen
.PHONY: gen
gen:
	go generate ./...
