#
# basic-app is intended to provide an example end-to-end use of hsc and hsc module in a node.js app
#


# One of: Darwin_i386, Darwin_x86_64, Linux_i386, Linux_x86_64
HSC_ARCH := Linux_x86_64
HSC_VERSION := 0.27.0
HSC_RELEASE_URL := "https://github.com/klyed/hivesmartchain/releases/download/v${HSC_VERSION}/hsc_${HSC_VERSION}_${HSC_ARCH}.tar.gz "
# Set to 'hsc' to use whatever is on PATH instead
HSC_BIN := bin/hsc

#
# Running the chain
#
# Make a simple single node chain
.PHONY: chain
chain: bin/hscburrow.toml

# Get the hsc binary
bin/hsc:
	mkdir -p bin
	curl -L ${HSC_RELEASE_URL} | tar zx -C bin hsc

# Generate the chain
hsc.toml genesis.json:
	${HSC_BIN} spec --full-accounts 1 | ${HSC_BIN} configure --genesis-spec=- --separate-genesis-doc=genesis.json > hsc.toml

# Dump account information to file for app
account.json: genesis.json
	jq  '.Accounts[] | select(.Name == "Full_0")' genesis.json > account.json

# Reset hsc state
.PHONY: reset_chain
reset_chain:
	rm -rf .hsc

# Remove hsc chain completely
.PHONY: remove_chain
remove_chain:
	rm -rf hsc.toml genesis.json .keys .hsc

# remake and reset chain
.PHONY: rechain
rechain: | remove_chain chain

.PHONY: start_chain
start_chain: chain
	${HSC_BIN} start -v0

.PHONY: restart
restart: | rechain start_chain

#
# Deploying the contract
#
deploy.output.json: simplestorage.sol deploy.yaml account.json
	${HSC_BIN} deploy --address $(shell jq '.Address' account.json) deploy.yaml

.PHONY: delete_deploy
delete_deploy:
	rm -rf deploy.output.json

.PHONY: deploy
deploy: deploy.output.json

.PHONY: redeploy
redeploy: | delete_deploy deploy.output.json

#
# Running the app
#

.PHONY: yarn_install
yarn_install:
	yarn install

.PHONY: start_app
start_app: yarn_install deploy
	node app.js
