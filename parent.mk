PARENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

include $(PARENT_DIR)/makefiles/service.mk
include $(PARENT_DIR)/makefiles/commons.mk
include $(PARENT_DIR)/makefiles/grpc.mk
include $(PARENT_DIR)/makefiles/mock.mk
include $(PARENT_DIR)/makefiles/lint.mk
