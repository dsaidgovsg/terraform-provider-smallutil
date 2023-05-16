go version  # go version go1.15.7 linux/amd6

tag="v0.4.3"
for osarch in \
  linux,amd64 \
  linux,386 \
  linux,arm64 \
  linux,arm \
  darwin,amd64 \
  darwin,arm64 \
  windows,amd64 \
  windows,386 \
  ; do
  goos="$(echo $osarch | cut -d ',' -f 1)"
  goarch="$(echo $osarch | cut -d ',' -f 2)"
  CGO_ENABLED=0 GOOS=${goos} GOARCH=${goarch} go build
  
  if [ "${goos}" = "windows" ]; then
    mv terraform-provider-smallutil.exe terraform-provider-smallutil_${tag}.exe
    zip terraform-provider-smallutil_${goos}_${goarch}_${tag}.zip terraform-provider-smallutil_${tag}.exe
    rm terraform-provider-smallutil_${tag}.exe
  else
    mv terraform-provider-smallutil terraform-provider-smallutil_${tag}
   zip terraform-provider-smallutil_${goos}_${goarch}_${tag}.zip terraform-provider-smallutil_${tag}
   rm terraform-provider-smallutil_${tag}
  fi
done
