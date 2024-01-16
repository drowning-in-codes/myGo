PROGRAM=GOimg
IMAEG_NAME=proanimer/goimg
all:
	@echo "The program is ${PROGRAM}" 
	@go fmt
	@make run

test:
	@go test 

run:
	@go run .

buildimg:
	-@docker rmi ${IMAEG_NAME} 
	@docker build -t ${IMAEG_NAME} .

pushimg:
	@docker push ${IMAEG_NAME}

vendor:
	@go vendor

clean:
	@go clean
.PHONY: all test run vendor buildimg clean pushimg


	