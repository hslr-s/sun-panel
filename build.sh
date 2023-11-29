#!/bin/bash

REPO=$(
  cd $(dirname $0)
  pwd
)
COMMIT_SHA=$(git rev-parse --short HEAD)
VERSION=$(git describe --tags)
# VERSION="0.1.1"
FRONTEND="false"
BINARY="false"
RELEASE="false"

debugInfo() {
  echo "Repo:           $REPO"
  echo "Build frontend:   $FRONTEND"
  echo "Build binary:   $BINARY"
  echo "Release:        $RELEASE"
  echo "Version:        $VERSION"
  echo "Commit:        $COMMIT_SHA"
}

buildFrontend() {
  cd $REPO
  pwd
#   npm install pnpm -g
  pnpm install
  pnpm run build
}

buildBackEndAssets() {
  cd $REPO/service
#   export PATH=$PATH:/root/go/bin
  go install -a -v github.com/go-bindata/go-bindata/...@latest
  go install -a -v github.com/elazarl/go-bindata-assetfs/...@latest
  go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/...
}

# buildBinary() {
#   cd $REPO/service
#   # mv "${REPO}/dist" "${REPO}/web"
#   go build -o "sun-panel" --ldflags="-X sun-panel/global.RUNCODE=release" main.go
# }

_build() {
  cd $REPO/service
  pwd 
  local osarch=$1
  IFS=/ read -r -a arr <<<"$osarch"
  os="${arr[0]}"
  arch="${arr[1]}"
  gcc="${arr[2]}"

  # Go build to build the binary.
  export GOOS=$os
  export GOARCH=$arch
  export CC=$gcc
  export CGO_ENABLED=1

  pathRelease=$REPO/release

  if [ -n "$VERSION" ]; then
    outPath="sun-panel_${VERSION}_${os}_${arch}"
  else
    outPath="sun-panel_${COMMIT_SHA}_${os}_${arch}"
  fi
  outname="${pathRelease}/${outPath}/sun-panel"
  go build -o "${outname}" --ldflags="-X sun-panel/global.RUNCODE=release" main.go

  cd $pathRelease
  # copy front file
  cp -r "${REPO}/dist" "${pathRelease}/${outPath}/web"

  echo "Release ${outPath}"
  if [ "$os" = "windows" ]; then
    mv $outname $outPath/sun-panel.exe
    zip -r "${pathRelease}/${outPath}.zip" $outPath
  else
    mv $outname $outPath/sun-panel
    tar -zcvf "${pathRelease}/${outPath}.tar.gz" $outPath
  fi
  rm -rf "${pathRelease}/${outPath}"
}

release() {
  cd $REPO/service
  ## List of architectures and OS to test coss compilation.
  SUPPORTED_OSARCH="linux/amd64/gcc linux/arm/arm-linux-gnueabihf-gcc windows/amd64/x86_64-w64-mingw32-gcc linux/arm64/aarch64-linux-gnu-gcc"

  echo "Release builds for OS/Arch/CC: ${SUPPORTED_OSARCH}"
  for each_osarch in ${SUPPORTED_OSARCH}; do
    _build "${each_osarch}"
  done
}

usage() {
  # echo "Usage: $0 [-f] [-c] [-b] [-r]" 1>&2
  echo "Usage: $0 [-f]  [-b] [-r]" 1>&2
  exit 1
}

while getopts "bfcrd" o; do
  case "${o}" in
  b)
    FRONTEND="true"
    BINARY="true"
    ;;
  f)
    FRONTEND="true"
    ;;
  c)
    BINARY="true"
    ;;
  r)
    FRONTEND="true"
    RELEASE="true"
    ;;
  d)
    DEBUG="true"
    ;;
  *)
    usage
    ;;
  esac
done
shift $((OPTIND - 1))

if [ "$DEBUG" = "true" ]; then
  debugInfo
fi

if [ "$FRONTEND" = "true" ]; then
  buildFrontend
fi

# if [ "$BINARY" = "true" ]; then
#   buildBinary
# fi

if [ "$RELEASE" = "true" ]; then
  buildBackEndAssets
  release
fi