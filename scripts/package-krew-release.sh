#!/usr/bin/env bash
# Build kubectl-cndro release archives for Krew (binary + LICENSE).
# Usage: VERSION=v0.1.0 ./scripts/package-krew-release.sh
# Artifacts go to dist/krew/<VERSION>/; run shasum -a 256 on them for krew/cndro.yaml.
set -euo pipefail

VERSION="${VERSION:-v0.1.0}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
OUT="${ROOT}/dist/krew/${VERSION}"
STAGE="${OUT}/.stage"

rm -rf "${OUT}"
mkdir -p "${OUT}"

build_one() {
	local goos="$1"
	local goarch="$2"
	local ext="$3"
	local binname="$4"
	local archive_name="$5"

	rm -rf "${STAGE}"
	mkdir -p "${STAGE}"

	CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" go build -trimpath -ldflags="-s -w" \
		-o "${STAGE}/${binname}" ./cmd/kubectl-cndro
	cp LICENSE "${STAGE}/"

	case "${ext}" in
	tar.gz)
		tar -czf "${OUT}/${archive_name}" -C "${STAGE}" "${binname}" LICENSE
		;;
	zip)
		( cd "${STAGE}" && zip -q "${OUT}/${archive_name}" "${binname}" LICENSE )
		;;
	esac
	echo "wrote ${OUT}/${archive_name}"
}

build_one linux amd64 tar.gz kubectl-cndro "kubectl-cndro_${VERSION}_linux_amd64.tar.gz"
build_one linux arm64 tar.gz kubectl-cndro "kubectl-cndro_${VERSION}_linux_arm64.tar.gz"
build_one darwin amd64 tar.gz kubectl-cndro "kubectl-cndro_${VERSION}_darwin_amd64.tar.gz"
build_one darwin arm64 tar.gz kubectl-cndro "kubectl-cndro_${VERSION}_darwin_arm64.tar.gz"
build_one windows amd64 zip kubectl-cndro.exe "kubectl-cndro_${VERSION}_windows_amd64.zip"

rm -rf "${STAGE}"

echo ""
echo "SHA-256 (paste into krew/cndro.yaml):"
shasum -a 256 "${OUT}"/*.tar.gz "${OUT}"/*.zip | sort -k2
