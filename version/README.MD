# version

编译时，记录程序版本相关信息

# 例子
```shell
version=$( git rev-parse --short HEAD 2> /dev/null || echo 'unknown' )
revision='alpha'
branch=$( git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown' )
BUILD_USER=${BUILD_USER:-"${USER}@${HOSTNAME}"}
BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}

build_info="-X github.com/sandwich-go/boost/version.BuildDate=$BUILD_DATE"
build_info="$build_info -X github.com/sandwich-go/boost/version.Revision=$revision"
build_info="$build_info -X github.com/sandwich-go/boost/version.Branch=$branch"
build_info="$build_info -X github.com/sandwich-go/boost/version.BuildUser=$BUILD_USER"
build_info="$build_info -X github.com/sandwich-go/boost/version.Version=$version"

go build -ldflags "$build_info" -o xxxxxx/main.go
```