#!/bin/bash

# pkg-config가 설치되어 있는지 확인
#if ! command -v pkg-config &> /dev/null; then
#    echo "pkg-config is not installed. Installing..."
#
#    # 업데이트 및 pkg-config 설치
#    sudo apt-get update
#    sudo apt-get install -y pkg-config
#
#    echo "pkg-config has been installed."
#else
#    echo "pkg-config is already installed."
#fi

# 서버 빌드 후 실행
go build -ldflags "-w -s" -o server
./server