#!/bin/bash

REPO=$(
  cd $(dirname $0)
  pwd
)
COMMIT_SHA=$(git rev-parse --short HEAD)
# VERSION=$(git describe --tags)
VERSION="v${cut -d '|' -f 2 ./service/assets/version}"
LATEST_TAG=$(git describe --tags --abbrev=0)
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
  echo "LATEST_TAG:        $LATEST_TAG"
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
  elif [ -n "$LATEST_TAG" ]; then
    outPath="sun-panel_${LATEST_TAG}_${os}_${arch}"
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

# 定义函数BuildReleaseLinuxMusl，用于构建正式版Linux-musl平台的二进制文件(参考Alist构建方案)
buildReleaseLinuxMusl() {
  cd $REPO/service
  ldflags="-X sun-panel/global.RUNCODE=release"
  pathRelease=$REPO/release
  # 清理.git目录，创建build目录，并下载交叉编译工具
  # rm -rf .git/
  # mkdir -p "build"
  muslflags="--extldflags '-static -fpic' $ldflags"
  BASE="https://musl.nn.ci/"
  # FILES=(x86_64-linux-musl-cross aarch64-linux-musl-cross mips-linux-musl-cross mips64-linux-musl-cross mips64el-linux-musl-cross mipsel-linux-musl-cross powerpc64le-linux-musl-cross s390x-linux-musl-cross)
  FILES=(x86_64-linux-musl-cross)
  for i in "${FILES[@]}"; do
    url="${BASE}${i}.tgz"
    curl -L -o "${i}.tgz" "${url}"
    tar xf "${i}.tgz" --strip-components 1 -C /usr/local
    rm -f "${i}.tgz"
  done
  # OS_ARCHES=(linux-musl-amd64 linux-musl-arm64 linux-musl-mips linux-musl-mips64 linux-musl-mips64le linux-musl-mipsle linux-musl-ppc64le linux-musl-s390x)
  # CGO_ARGS=(x86_64-linux-musl-gcc aarch64-linux-musl-gcc mips-linux-musl-gcc mips64-linux-musl-gcc mips64el-linux-musl-gcc mipsel-linux-musl-gcc powerpc64le-linux-musl-gcc s390x-linux-musl-gcc)

  # 暂时仅编译amd64
  OS_ARCHES=(linux-musl-amd64)
  CGO_ARGS=(x86_64-linux-musl-gcc)
  

  for i in "${!OS_ARCHES[@]}"; do
    os_arch=${OS_ARCHES[$i]}
    cgo_cc=${CGO_ARGS[$i]}
    echo building for ${os_arch}
    export GOOS=${os_arch%%-*}
    export GOARCH=${os_arch##*-}
    export CC=${cgo_cc}
    export CGO_ENABLED=1

    if [ -n "$VERSION" ]; then
      outPath="sun-panel_${VERSION}_${GOOS}_musl_${GOARCH}"
    elif [ -n "$LATEST_TAG" ]; then
      outPath="sun-panel_${LATEST_TAG}_${GOOS}_musl_${GOARCH}"
    else
      outPath="sun-panel_${COMMIT_SHA}_${GOOS}_musl_${GOARCH}"
    fi

    outname="${pathRelease}/${outPath}/sun-panel"

    go build -o "${outname}" -ldflags="$muslflags" main.go
    # go build -o "${outname}" -ldflags="$muslflags" -tags=jsoniter main.go
  done

  cd $pathRelease
  # copy front file
  cp -r "${REPO}/dist" "${pathRelease}/${outPath}/web"

  echo "Release ${outPath}"

  mv $outname $outPath/sun-panel
  tar -zcvf "${pathRelease}/${outPath}.tar.gz" $outPath
  
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

  # 临时方案解决centos无法运行的问题
  buildReleaseLinuxMusl
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