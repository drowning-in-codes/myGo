PROGRAM=GOimg

all:
	@echo "The program is ${PROGRAM}" 
	@go fmt
	@make run

test:
	@go test 

run:
	@go run .

vendor:
	@go vendor

.PHONY: all test run vendor


	