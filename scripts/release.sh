#!/usr/bin/env bash
# Builds a release into OUTDIR: an archive per platform holding gopamin,
# LICENSE, NOTICES.txt and a SHA256SUMS of those files; the third-party
# notices; the rendered install.sh; and a SHA256SUMS of all of them. The
# release workflow runs it for a tag; it runs the same way locally.
set -euo pipefail

usage() {
	echo "usage: scripts/release.sh VERSION [OUTDIR]" >&2
	echo "  VERSION like 1.3.0, equal to VERSION in internal/commands/constants.go without the v; OUTDIR defaults to dist" >&2
	exit 2
}

[ $# -ge 1 ] && [ $# -le 2 ] || usage
version=$1
out=${2:-dist}
[[ $version =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]] || { echo "release: invalid version $version" >&2; usage; }

root=$(cd "$(dirname "$0")/.." && pwd)
cd "$root"

declared=$(sed -n 's/^const VERSION = "v\([^"]*\)"$/\1/p' internal/commands/constants.go)
if [ "$declared" != "$version" ]; then
	echo "release: internal/commands/constants.go declares VERSION v$declared, not v$version; bump it first" >&2
	exit 1
fi

if [ -n "$(git status --porcelain)" ]; then
	echo "release: the working tree has uncommitted changes; commit or stash them first" >&2
	exit 1
fi
commit=$(git rev-parse HEAD)

if [ -e "$out" ]; then
	echo "release: $out already exists; remove it or choose another directory" >&2
	exit 1
fi
mkdir -p "$out"
out=$(cd "$out" && pwd)

if command -v sha256sum >/dev/null 2>&1; then
	sum() { sha256sum "$@"; }
else
	sum() { shasum -a 256 "$@"; }
fi

targets="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"
export GOWORK=off CGO_ENABLED=0

goroot=$(go env GOROOT)
golicense=""
for f in "$goroot/LICENSE" "$goroot/../LICENSE"; do
	if [ -f "$f" ]; then
		golicense=$f
		break
	fi
done
if [ -z "$golicense" ]; then
	echo "release: the Go toolchain's LICENSE wasn't found near $goroot" >&2
	exit 1
fi

notices="$out/gopamin-$version-NOTICES.txt"
{
	echo "Gopamin v$version, built from commit $commit with $(go version | cut -d' ' -f3)"
	echo
	echo "==== Gopamin ===="
	cat LICENSE
	echo
	echo "==== Go standard library ($(go version | cut -d' ' -f3)) ===="
	cat "$golicense"
	for t in $targets; do
		GOOS=${t%/*} GOARCH=${t#*/} go list -deps -f '{{with .Module}}{{if not .Main}}{{.Path}} {{.Version}} {{.Dir}}{{end}}{{end}}' .
	done | sort -u |
		while read -r path ver dir; do
			found=0
			for f in "$dir"/LICENSE* "$dir"/NOTICE* "$dir"/COPYING*; do
				[ -f "$f" ] || continue
				echo
				echo "==== $path $ver: $(basename "$f") ===="
				cat "$f"
				found=1
			done
			if [ $found -eq 0 ]; then
				echo "release: no license file found for $path $ver" >&2
				exit 1
			fi
		done
} > "$notices"
echo "wrote $notices"

stage="$out/.stage"
trap 'rm -rf "$stage"' EXIT

for t in $targets; do
	os=${t%/*}
	arch=${t#*/}
	name="gopamin-$version-$os-$arch"
	pkg="$stage/$name"
	mkdir -p "$pkg"
	GOOS=$os GOARCH=$arch go build -trimpath -buildvcs=false -ldflags "-s -w -buildid=" -o "$pkg/gopamin" .
	cp LICENSE "$pkg/LICENSE"
	cp "$notices" "$pkg/NOTICES.txt"
	(cd "$pkg" && sum gopamin LICENSE NOTICES.txt > "$stage/SHA256SUMS" && mv "$stage/SHA256SUMS" SHA256SUMS)
	COPYFILE_DISABLE=1 tar -C "$stage" -czf "$out/$name.tar.gz" "$name"
	echo "built $out/$name.tar.gz"
done
rm -rf "$stage"
trap - EXIT

(cd "$out" && sum gopamin-"$version"-*.tar.gz > SHA256SUMS)
sh "$root/scripts/release-install.sh" "$version" "$out"
(cd "$out" && sum gopamin-"$version"-* install.sh > SHA256SUMS.part && mv SHA256SUMS.part SHA256SUMS)
echo "wrote $out/SHA256SUMS"
cat "$out/SHA256SUMS"
