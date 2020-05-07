#!/bin/bash

platforms=(linux-386 linux-amd64 linux-arm linux-arm64 darwin-amd64 windows-amd64 windows-386)

function build {
  if [ "$1" == "windows" ]; then
    local suffix=".exe"
  fi

  GOOS=$1 GOARCH=$2 go build -o bin/getty$suffix -ldflags="-s -w" getty.go
  
  if [ "$1" == "windows" ]; then
    zip -j -q pkg/getty-"$1"-"$2".zip bin/getty$suffix
  else
    tar czf pkg/getty-"$1"-"$2".tar.gz -C bin getty
  fi
}

rm -rf bin/*
rm -rf pkg/*
mkdir -p bin pkg

for i in "${platforms[@]}"
do
  IFS='-'
  read -a strarr <<< "$i"
  build "${strarr[0]}" "${strarr[1]}"
done

rm -r bin