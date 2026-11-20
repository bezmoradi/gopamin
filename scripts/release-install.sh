#!/bin/sh
# Renders scripts/install.sh for a release: puts the version and the four
# archives' SHA-256 sums from DIST/SHA256SUMS in place of the placeholders.
set -eu
LC_ALL=C
export LC_ALL
template=$(dirname -- "$0")/install.sh
usage='usage: scripts/release-install.sh VERSION DIST (DIST holds the archives and their SHA256SUMS; writes DIST/install.sh)'
[ $# = 2 ] || {
	echo "$usage" >&2
	exit 2
}
version=$1
dist=$2
printf '%s\n' "$version" | grep -Eqx '[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?' || {
	echo "release-install: not a version MAJOR.MINOR.PATCH or MAJOR.MINOR.PATCH-PRERELEASE: $version" >&2
	exit 1
}
sums="$dist/SHA256SUMS"
[ -s "$sums" ] || {
	echo "release-install: no such file, or empty: $sums" >&2
	exit 1
}
sha256() {
	if command -v sha256sum >/dev/null 2>&1; then
		sum=$(sha256sum "$1")
	else
		sum=$(shasum -a 256 "$1")
	fi
	printf '%s\n' "${sum%% *}"
}
script="s|^version='@VERSION@'\$|version='$version'|"
check=
for target in linux-amd64 linux-arm64 darwin-amd64 darwin-arm64; do
	name="gopamin-$version-$target.tar.gz"
	key=$(printf '%s' "$target" | tr - _)
	placeholder=$(printf '%s' "$key" | tr '[:lower:]' '[:upper:]')
	listed=$(awk -v n="$name" 'NF == 2 && $2 == n { print $1 }' "$sums")
	printf '%s\n' "$listed" | grep -Eqx '[0-9a-f]{64}' || {
		echo "release-install: $sums doesn't list one SHA-256 for $name" >&2
		exit 1
	}
	if [ ! -f "$dist/$name" ] || [ "$(sha256 "$dist/$name")" != "$listed" ]; then
		echo "release-install: $dist/$name is missing or doesn't match its SHA-256 in $sums" >&2
		exit 1
	fi
	[ "$(grep -cxF "sha256_$key='@SHA256_$placeholder@'" "$template")" = 1 ] || {
		echo "release-install: $template doesn't hold the line sha256_$key='@SHA256_$placeholder@' exactly once" >&2
		exit 1
	}
	script="$script
s|^sha256_$key='@SHA256_$placeholder@'\$|sha256_$key='$listed'|"
	check="$check sha256_$key='$listed'"
done
[ "$(grep -cxF "version='@VERSION@'" "$template")" = 1 ] || {
	echo "release-install: $template doesn't hold the line version='@VERSION@' exactly once" >&2
	exit 1
}
out="$dist/install.sh"
part="$out.part"
sed -e "$script" "$template" >"$part"
ok=1
grep -q '@[A-Z0-9_]*@' "$part" && ok=0
[ "$(grep -cxF "version='$version'" "$part")" = 1 ] || ok=0
for line in $check; do
	[ "$(grep -cxF "$line" "$part")" = 1 ] || ok=0
done
sh -n "$part" || ok=0
if [ "$ok" != 1 ]; then
	echo "release-install: $part doesn't carry the version and the four SHA-256 sums in place of the placeholders, or doesn't parse" >&2
	rm -f "$part"
	exit 1
fi
mv -f "$part" "$out"
echo "release-install: wrote $out for $version"
grep -E "^(version|sha256_[a-z0-9_]+)='" "$out"
