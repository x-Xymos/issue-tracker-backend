BUILD_DIR=build/

SRC_NAME1=./src/services/account-api
BINARY_NAME1=account-api

SRC_NAME2=./src/services/project-api
BINARY_NAME2=project-api

build:
	go build -o  $(BUILD_DIR)$(BINARY_NAME1) -v $(SRC_NAME1)
	go build -o  $(BUILD_DIR)$(BINARY_NAME2) -v $(SRC_NAME2)

run:
	nohup bash ./run_service $(BUILD_DIR) $(BINARY_NAME1) &
	nohup bash ./run_service $(BUILD_DIR) $(BINARY_NAME2) &


