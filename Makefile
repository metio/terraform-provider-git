# SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
# SPDX-License-Identifier: 0BSD

###############################################################################
# PROLOGUE                                                                    #
###############################################################################
MAKEFLAGS += --warn-undefined-variables
.DEFAULT_GOAL := help

###############################################################################
# HELP                                                                        #
###############################################################################

##@ other
.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

###############################################################################
# PROJECT SPECIFIC                                                            #
###############################################################################
-include project.mk
