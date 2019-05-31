set -euo pipefail

REPO_OWNER=guangie88
REPO_NAME=terraform-provider-smallutil
REPO=${REPO_OWNER}/${REPO_NAME}

# https://gist.github.com/lukechilds/a83e1d7127b78fef38c2914c4ececc3c
get_latest_release() {
  curl --silent "https://api.github.com/repos/${REPO}/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                                 # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                         # Pluck JSON value
}

TAG=$(get_latest_release)
echo "Latest tag is: ${TAG}"

# Detect the OS and ARCH for download
OS_CALL="$(uname -s)"
if [ "${OS_CALL}" = "Linux" ]; then
  OS=linux
elif [ "${OS_CALL}" == "Darwin" ]; then
  OS=darwin
else
  echo "Unable to continue for OS type ${OS_CALL}..."
  exit 1
fi

ARCH_CALL="$(uname -m)"
if [[ "${ARCH_CALL}" == i*86 ]]; then
  ARCH=386
elif [[ "${ARCH_CALL}" == x86_64* ]]; then
  ARCH=amd64
elif [[ "${ARCH_CALL}" == arm*64 ]]; then
  ARCH=arm64
elif [[ "${ARCH_CALL}" == arm* ]]; then
  ARCH=arm
else
  echo "Unable to continue for ARCH type ${ARCH_CALL}..."
  exit 1
fi

# Download archive containing plugin
URL=https://github.com/${REPO}/releases/download/${TAG}/${REPO_NAME}_${OS}_${ARCH}_${TAG}.zip
echo "Downloading from ${URL}..."
curl -sLO ${URL}
echo "Download complete!"

# Unzip and move it into Terraform plugin dir
unzip ${REPO_NAME}_${OS}_${ARCH}_${TAG}.zip
mkdir -p ${HOME}/.terraform.d/plugins
mv ${REPO_NAME}_${TAG} ${HOME}/.terraform.d/plugins/
rm ${REPO_NAME}_${OS}_${ARCH}_${TAG}.zip

echo "Plugin installation completed!"
