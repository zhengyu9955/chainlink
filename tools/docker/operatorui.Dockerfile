# Build Chainlink
FROM smartcontract/builder:1.0.27-rust-1.10-test

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}
