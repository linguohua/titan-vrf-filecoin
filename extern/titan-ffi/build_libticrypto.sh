#!/usr/bin/env bash
set -e

echo -e "\033[34m build libticrypto.a \033[0m"
cd rust && \
./scripts/build-release.sh build && \
#cargo build --package ticrypto --lib --release && \
cp target/release/libticrypto.a ../libticrypto.a && \
cp ticrypto.pc ../ticrypto.pc && \
cp ticrypto.h ../ticrypto.h
