#!/bin/bash
projRoot=$(realpath $1)

go_module=github.com/shenrytech/shdb

# Find location of well-known types
well_known_types_path=$(dirname $(which protoc))/../include

function build_pb_go() {
    local file=$1
    printf "golang"
    protoc \
        -I $well_known_types_path \
        -I $projRoot \
        --go_opt=module=$go_module \
        --go_out=$projRoot \
        --go-grpc_out=$projRoot \
        --go-grpc_opt=module=$go_module \
        $file

}

function build_pb() {
    local file=$1
    build_pb_go $1
    printf " "
}

build_pb ${projRoot}/shdb.proto
