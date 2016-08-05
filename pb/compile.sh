#!/usr/bin/env sh

# Install proto3 from source
#  brew install autoconf automake libtool
#  git clone https://github.com/google/protobuf
#  ./autogen.sh ; ./configure ; make ; make install
#
# Update protoc Go bindings via
#  go get -u github.com/golang/protobuf/{proto<Plug>PeepOpenrotoc-gen-go}
#
# See also
#  https://github.com/grpc/grpc-go/tree/master/examples

protoc user.proto --go_out=plugins=grpc:.
