#!/bin/sh
# Installs Gopamin from a GitHub release. This file is a template: the release
# workflow renders it with the version and the archives' SHA-256 sums
# (scripts/release-install.sh) and attaches the result to the release.
set -eu

version='@VERSION@'
sha256_linux_amd64='@SHA256_LINUX_AMD64@'
sha256_linux_arm64='@SHA256_LINUX_ARM64@'
sha256_darwin_amd64='@SHA256_DARWIN_AMD64@'
sha256_darwin_arm64='@SHA256_DARWIN_ARM64@'

repo=bezmoradi/gopamin
latest="https://github.com/$repo/releases/latest/download/install.sh"
base="https://github.com/$repo/releases/download/v$version"
go_needed=1.27.3
go_switching=1.21.0
LC_ALL=C
export LC_ALL

tmp=
stage=

stop() {
	printf 'Install stopped: %s\n' "$1" >&2
	exit 1
}

cleanup() {
	[ -z "$stage" ] || rm -rf "$stage"
	[ -z "$tmp" ] || rm -rf "$tmp"
}

sha256_of() {
	if command -v sha256sum >/dev/null 2>&1; then
		sum=$(sha256sum "$1") || return 1
	else
		sum=$(shasum -a 256 "$1") || return 1
	fi
	printf '%s\n' "${sum%% *}"
}

check_sums() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum -c SHA256SUMS >/dev/null 2>&1
	else
		shasum -a 256 -c SHA256SUMS >/dev/null 2>&1
	fi
}

# go_at_least GOVERSION MAJOR.MINOR.PATCH succeeds when a GOVERSION such as
# go1.27.3, go1.27.10, go1.28rc1 (counted as 1.28.0) or "devel go1.28-abc" is
# that version or later.
go_at_least() {
	gv_want=$2
	gv_parts=$(printf '%s\n' "$1" | sed -n 's/^[^0-9]*go\([0-9][0-9]*\)\.\([0-9][0-9]*\)\(\.\([0-9][0-9]*\)\)\{0,1\}.*$/\1 \2 \4/p' | head -n 1)
	[ -n "$gv_parts" ] || return 1
	# shellcheck disable=SC2086
	set -- $gv_parts
	gv_major=$1 gv_minor=$2 gv_patch=${3:-0}
	[ "${#gv_major}" -le 6 ] && [ "${#gv_minor}" -le 6 ] && [ "${#gv_patch}" -le 6 ] || return 1
	gv_want_major=${gv_want%%.*}
	gv_want_patch=${gv_want##*.}
	gv_want_minor=${gv_want#*.}
	gv_want_minor=${gv_want_minor%.*}
	[ "$gv_major" -ne "$gv_want_major" ] && { [ "$gv_major" -gt "$gv_want_major" ]; return; }
	[ "$gv_minor" -ne "$gv_want_minor" ] && { [ "$gv_minor" -gt "$gv_want_minor" ]; return; }
	[ "$gv_patch" -ge "$gv_want_patch" ]
}

main() {
	exec </dev/null
	case $version$sha256_linux_amd64$sha256_linux_arm64$sha256_darwin_amd64$sha256_darwin_arm64 in
	*@*) stop "this is the template in the repository, not a release's script. Install with: curl -fsSL $latest | sh" ;;
	esac

	system=$(uname -s)
	machine=$(uname -m)
	case $system in
	Linux) os=linux ;;
	Darwin) os=darwin ;;
	*) os= ;;
	esac
	case $machine in
	x86_64 | amd64) arch=amd64 ;;
	arm64 | aarch64) arch=arm64 ;;
	*) arch= ;;
	esac
	if [ -z "$os" ] || [ -z "$arch" ]; then
		stop "this is $system on $machine, and Gopamin's releases are built only for linux/amd64, linux/arm64, darwin/amd64 and darwin/arm64."
	fi
	case $os-$arch in
	linux-amd64) want=$sha256_linux_amd64 ;;
	linux-arm64) want=$sha256_linux_arm64 ;;
	darwin-amd64) want=$sha256_darwin_amd64 ;;
	darwin-arm64) want=$sha256_darwin_arm64 ;;
	esac

	command -v curl >/dev/null 2>&1 || stop "curl isn't installed. Install curl and run this again."
	command -v tar >/dev/null 2>&1 || stop "tar isn't installed. Install tar and run this again."
	if ! command -v sha256sum >/dev/null 2>&1 && ! command -v shasum >/dev/null 2>&1; then
		stop "neither sha256sum nor shasum is installed, so the download can't be checked. Install coreutils (or perl's shasum) and run this again."
	fi
	command -v git >/dev/null 2>&1 || stop "git isn't installed, and Gopamin starts a Git repository in every project it creates. Install git and run this again."
	command -v go >/dev/null 2>&1 || stop "Go isn't installed, and Gopamin runs it to create projects. Install Go $go_needed or later from https://go.dev/dl and run this again."
	goversion=$(cd / && go env GOVERSION 2>/dev/null | head -n 1) || goversion=
	gotoolchain=$(cd / && go env GOTOOLCHAIN 2>/dev/null | head -n 1) || gotoolchain=
	go_note=
	if ! go_at_least "$goversion" "$go_needed"; then
		if ! go_at_least "$goversion" "$go_switching"; then
			stop "Go here reports version '$goversion', and Gopamin needs Go 1.21 or later, which fetches Go $go_needed by itself, to create projects. Install it from https://go.dev/dl and run this again."
		fi
		case $gotoolchain in
		auto | *+auto) go_note="Go here is $goversion, so it will download go$go_needed the first time Gopamin creates a project." ;;
		*) stop "Go here is $goversion, and GOTOOLCHAIN=$gotoolchain keeps it from fetching go$go_needed, which Gopamin needs to create projects. Install Go $go_needed or later from https://go.dev/dl, or set GOTOOLCHAIN=auto, and run this again." ;;
		esac
	fi

	if [ -n "${GOPAMIN_INSTALL_DIR:-}" ]; then
		dir=$GOPAMIN_INSTALL_DIR
	elif [ -n "${HOME:-}" ]; then
		dir=$HOME/.local/bin
	else
		stop "HOME isn't set. Set GOPAMIN_INSTALL_DIR to the folder to install into."
	fi
	case $dir in
	/*) ;;
	*) dir="$(pwd)/$dir" ;;
	esac

	trap cleanup EXIT
	trap 'exit 1' HUP INT TERM
	tmp=$(mktemp -d "${TMPDIR:-/tmp}/gopamin-install.XXXXXX" 2>/dev/null) || stop "couldn't make a temporary folder."

	name="gopamin-$version-$os-$arch"
	archive="$tmp/$name.tar.gz"
	echo "Downloading $base/$name.tar.gz"
	curl -fsSL --proto '=https' --tlsv1.2 --retry 3 -o "$archive" "$base/$name.tar.gz" || stop "couldn't download $base/$name.tar.gz; nothing in $dir was changed."
	if ! got=$(sha256_of "$archive") || [ "$got" != "$want" ]; then
		stop "$name.tar.gz doesn't match the SHA-256 this script carries, so nothing in $dir was changed."
	fi
	mkdir "$tmp/x" || stop "couldn't make a folder in $tmp."
	tar -xzf "$archive" -C "$tmp/x" --no-same-owner || stop "couldn't unpack $name.tar.gz; nothing in $dir was changed."
	pkg="$tmp/x/$name"
	[ -f "$pkg/SHA256SUMS" ] || stop "$name.tar.gz holds no $name/SHA256SUMS; nothing in $dir was changed."
	if grep -Evx '[0-9a-f]{64}  [A-Za-z0-9._-]+' "$pkg/SHA256SUMS" >/dev/null; then
		stop "$name/SHA256SUMS has a line that isn't a SHA-256 and a file name; nothing in $dir was changed."
	fi
	(cd "$pkg" && check_sums) || stop "the files in $name.tar.gz don't match its SHA256SUMS; nothing in $dir was changed."
	grep -Eqx '[0-9a-f]{64}  gopamin' "$pkg/SHA256SUMS" || stop "$name/SHA256SUMS doesn't list gopamin; nothing in $dir was changed."
	if [ -L "$pkg/gopamin" ] || [ ! -f "$pkg/gopamin" ] || [ ! -x "$pkg/gopamin" ]; then
		stop "gopamin in $name.tar.gz isn't a program; nothing in $dir was changed."
	fi
	got=$("$pkg/gopamin" version 2>/dev/null) || stop "the downloaded gopamin doesn't run from $tmp; if TMPDIR doesn't allow running programs, set it to a folder that does. Nothing in $dir was changed."
	[ "$got" = "Gopamin CLI Version: v$version" ] || stop "the downloaded gopamin reports '$got', not version v$version; nothing in $dir was changed."

	mkdir -p "$dir" 2>/dev/null || stop "couldn't make the folder $dir."
	if [ -d "$dir/gopamin" ]; then
		stop "$dir/gopamin is a folder or a link to one; move it away and run this again. Nothing in $dir was changed."
	fi
	stage=$(mktemp -d "$dir/.gopamin-install.XXXXXX" 2>/dev/null) || stop "couldn't write in $dir; nothing in it was changed."
	if ! cp "$pkg/gopamin" "$stage/gopamin" || ! chmod 0755 "$stage/gopamin"; then
		stop "couldn't copy gopamin into $dir; nothing in it was changed."
	fi
	mv -f "$stage/gopamin" "$dir/gopamin" || stop "couldn't move gopamin into place in $dir; nothing in it was changed."
	rm -rf "$stage"
	stage=

	echo "Installed Gopamin v$version for $os/$arch in $dir/gopamin"
	installed=$("$dir/gopamin" version) || stop "$dir/gopamin doesn't run."
	echo "$dir/gopamin version: $installed"
	[ -z "$go_note" ] || echo "Note: $go_note"
	case ":${PATH:-}:" in
	*":$dir:"*) ;;
	*) echo "$dir isn't on PATH; add it in your shell's profile: export PATH=\"$dir:\$PATH\"" ;;
	esac
	echo "To check where this release came from, with a signed-in gh: curl -fsSLO $base/$name.tar.gz && gh attestation verify $name.tar.gz -R $repo --signer-workflow $repo/.github/workflows/release.yml --source-ref refs/tags/v$version"
}

main "$@"
