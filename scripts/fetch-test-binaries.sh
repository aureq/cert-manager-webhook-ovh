#!/bin/bash

set -e

k8s_version=1.32.0
arch=amd64

if [ "$OSTYPE" = "linux-gnu" ]; then
  os="linux"
elif [ "$OSTYPE" = "darwin"* ]; then
  os="darwin"
else
  echo "OS '$OSTYPE' not supported." >&2
  exit 1
fi

root=$(cd "`dirname $0`"/..; pwd)
output_dir="$root"/_out
archive_name="kubebuilder-tools-$k8s_version-$os-$arch.tar.gz"
archive_file="$output_dir/$archive_name"
archive_url="https://storage.googleapis.com/kubebuilder-tools/$archive_name"

mkdir -p "$output_dir"
curl -sL "$archive_url" -o "$archive_file"
tar -zxf "$archive_file" -C "$output_dir/"
