#!/bin/bash
# Sourced by the offline generators before they rm -rf a directory.

# Absolute path with every existing component resolved through symlinks.
# Components that do not exist yet are appended unchanged.
physical_path() {
  local existing="$1" missing=""
  while [[ "$existing" == ?*/ ]]; do
    existing="${existing%/}"
  done
  [[ "$existing" == /* ]] || existing="$PWD/$existing"
  while [[ ! -d "$existing" ]]; do
    missing="/${existing##*/}$missing"
    existing="${existing%/*}"
    existing="${existing:-/}"
  done
  existing="$(cd -- "$existing" && pwd -P)" || return
  printf '%s%s\n' "${existing%/}" "$missing"
}

# True when PATH, or a directory above it, is the existing directory DIR.
# -ef compares the directories themselves, so another spelling of DIR on a
# case-insensitive volume (Offline) is still DIR.
at_or_under() {
  local path="$1"
  while :; do
    [[ "$path" -ef "$2" ]] && return 0
    [[ "$path" == / ]] && return 1
    path="${path%/*}"
    path="${path:-/}"
  done
}

# require_deletable_dir NAME VALUE [full]
# Exits unless what rm -rf VALUE removes lies strictly inside this scripts
# directory and clear of the full host tile set, scripts/offline. A full
# generator passes "full" to replace that set as a whole.
require_deletable_dir() {
  local name="$1" value="$2" scripts entry tiles target parent id link=0
  scripts="$(physical_path "$(dirname -- "${BASH_SOURCE[0]}")")"
  entry="$scripts/offline"
  tiles="$(physical_path "$entry")"
  # rm -rf removes a symlink named without a trailing slash, not its target.
  if [[ "$value" != */ && -L "$value" ]]; then
    target="$(physical_path "$(dirname -- "$value")")/${value##*/}"
    link=1
  else
    target="$(physical_path "$value")"
  fi
  parent="${target%/*}"
  parent="${parent:-/}"
  case "/$value/" in
    */../*) ;;
    *)
      if at_or_under "$parent" "$scripts"; then
        # A link is the tile set's own entry when it is that same link, whatever its spelling.
        if [[ "$target" == "$entry" ]] || { (( ! link )) && [[ "$target" -ef "$entry" ]]; } ||
          { (( link )) && [[ -L "$entry" ]] && id="$(/usr/bin/stat -f %d:%i -- "$value")" &&
            [[ -n "$id" && "$id" == "$(/usr/bin/stat -f %d:%i -- "$entry")" ]]; }; then
          [[ "${3:-}" == full ]] && return 0
        elif (( link )); then
          at_or_under "$parent" "$tiles" || return 0
        elif ! at_or_under "$target" "$tiles" && ! at_or_under "$tiles" "$target"; then
          return 0
        fi
      fi
      ;;
  esac
  echo "refusing unsafe $name: $value (removes $target)" >&2
  exit 1
}
