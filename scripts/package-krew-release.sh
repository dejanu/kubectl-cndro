#!/usr/bin/env bash
# Build kubectl-cndro release archives for Krew (binary + LICENSE).
# Usage: VERSION=v0.1.1 ./scripts/package-krew-release.sh
# Artifacts go to dist/krew/<VERSION>/; manifest plugins/cndro.yaml
# is updated automatically with the new version and sha256 values.
set -euo pipefail

# sha256 helper: works on both Linux (sha256sum) and macOS (shasum)
_sha256() {
	if command -v sha256sum &>/dev/null; then
		sha256sum "$1" | awk '{print $1}'
	else
		shasum -a 256 "$1" | awk '{print $1}'
	fi
}

VERSION="${VERSION:-v0.1.1}"
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

	CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" \
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

# Capture sha256 values per platform
SHA256_LINUX_AMD64="$(_sha256 "${OUT}/kubectl-cndro_${VERSION}_linux_amd64.tar.gz")"
SHA256_LINUX_ARM64="$(_sha256 "${OUT}/kubectl-cndro_${VERSION}_linux_arm64.tar.gz")"
SHA256_DARWIN_AMD64="$(_sha256 "${OUT}/kubectl-cndro_${VERSION}_darwin_amd64.tar.gz")"
SHA256_DARWIN_ARM64="$(_sha256 "${OUT}/kubectl-cndro_${VERSION}_darwin_arm64.tar.gz")"
SHA256_WINDOWS_AMD64="$(_sha256 "${OUT}/kubectl-cndro_${VERSION}_windows_amd64.zip")"

echo ""
echo "SHA-256:"
echo "  linux_amd64:   ${SHA256_LINUX_AMD64}"
echo "  linux_arm64:   ${SHA256_LINUX_ARM64}"
echo "  darwin_amd64:  ${SHA256_DARWIN_AMD64}"
echo "  darwin_arm64:  ${SHA256_DARWIN_ARM64}"
echo "  windows_amd64: ${SHA256_WINDOWS_AMD64}"

# Update version and sha256 values in a Krew manifest file (YAML) in-place.
update_manifest() {
	python3 - "$1" "${VERSION}" \
		"${SHA256_LINUX_AMD64}" \
		"${SHA256_LINUX_ARM64}" \
		"${SHA256_DARWIN_AMD64}" \
		"${SHA256_DARWIN_ARM64}" \
		"${SHA256_WINDOWS_AMD64}" <<'PYEOF'
import sys, re

filename, version, *shas = sys.argv[1:]
platform_order = ['linux_amd64', 'linux_arm64', 'darwin_amd64', 'darwin_arm64', 'windows_amd64']
platform_shas = dict(zip(platform_order, shas))

with open(filename) as f:
    lines = f.readlines()

current_platform = None
out = []
for line in lines:
    # Update top-level version field
    if re.match(r'\s*version:\s*v', line):
        line = re.sub(r'(version:\s*)v\S+', r'\g<1>' + version, line)
    # Detect platform from uri line and rewrite version in URL + archive name
    m = re.search(r'uri:.*kubectl-cndro_v[^_]+_(\w+)_(\w+)\.(tar\.gz|zip)', line)
    if m:
        current_platform = m.group(1) + '_' + m.group(2)
        line = re.sub(r'(releases/download/)v[^/]+/', r'\g<1>' + version + '/', line)
        line = re.sub(r'(kubectl-cndro_)v[^_]+(_\w+_\w+\.)', r'\g<1>' + version + r'\2', line)
    # Update sha256 for the current platform
    elif re.match(r'\s*sha256:', line) and current_platform in platform_shas:
        line = re.sub(r'(sha256:\s*)\S+', r'\g<1>' + platform_shas[current_platform], line)
    out.append(line)

with open(filename, 'w') as f:
    f.writelines(out)
print('updated', filename)
PYEOF
}

echo ""
update_manifest "${ROOT}/plugins/cndro.yaml"
