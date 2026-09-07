#!/bin/bash
set -euo pipefail

MODE="$1"
OVERRIDE="$2"
GITHUB_SHA="$3"
ORIGIN_MAIN_SHA="$4"

if [[ "${GITHUB_SHA}" != "${ORIGIN_MAIN_SHA}" ]]; then
  echo "Error: Validation failed. ${GITHUB_SHA} is not the current origin/main (${ORIGIN_MAIN_SHA})" >&2
  exit 1
fi

VALID_TAG_REGEX='^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$'

if [[ -n "$OVERRIDE" ]]; then
  OVERRIDE="${OVERRIDE#v}"
  next_tag="v$OVERRIDE"

  if [[ ! "$next_tag" =~ $VALID_TAG_REGEX ]]; then
      echo "Error: Override tag does not meet required shape: $next_tag" >&2
      exit 1
  fi

  if git rev-parse "$next_tag" >/dev/null 2>&1; then
    REMOTE_TAG_SHA=$(git ls-remote --tags origin "refs/tags/$next_tag^{}" | awk '{print $1}')
    if [[ -z "$REMOTE_TAG_SHA" ]]; then
        REMOTE_TAG_SHA=$(git ls-remote --tags origin "refs/tags/$next_tag" | awk '{print $1}')
    fi

    if [[ "$REMOTE_TAG_SHA" == "${GITHUB_SHA}" ]]; then
      echo "Recovery Tag $next_tag exists and points to GITHUB_SHA. Safe to retry." >&2
      echo "$next_tag"
      exit 0
    else
      echo "Recovery Tag already exists and points to $REMOTE_TAG_SHA (expected ${GITHUB_SHA}): $next_tag" >&2
      exit 1
    fi
  else
    echo "Error: Recovery tag $next_tag does not exist. The override is for recovery only." >&2
    exit 1
  fi

else
  # Find latest stable tag
  latest_stable=$(git tag -l "v[0-9]*.[0-9]*.[0-9]*" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -n 1 || true)
  if [[ -z "$latest_stable" ]]; then latest_stable="v0.0.0"; fi

  clean_stable="${latest_stable#v}"
  IFS='.' read -r s_major s_minor s_patch <<< "$clean_stable"
  s_major=${s_major:-0}; s_minor=${s_minor:-0}; s_patch=${s_patch:-0}

  case "$MODE" in
    release-major) next_tag="v$((s_major + 1)).0.0" ;;
    release-minor) next_tag="v${s_major}.$((s_minor + 1)).0" ;;
    release-patch) next_tag="v${s_major}.${s_minor}.$((s_patch + 1))" ;;
    release-rc)
      next_stable_patch="v${s_major}.${s_minor}.$((s_patch + 1))"
      # Find specific latest RC for the target version
      latest_rc=$(git tag -l "${next_stable_patch}-rc*" | grep -E "^${next_stable_patch}-rc[0-9]+$" | sort -V | tail -n 1 || true)

      if [[ -n "$latest_rc" ]]; then
          rc_num=$(echo "$latest_rc" | grep -oE "rc[0-9]+" | sed 's/rc//')
          next_tag="${next_stable_patch}-rc$((rc_num + 1))"
      else
          next_tag="${next_stable_patch}-rc1"
      fi
      ;;
    release-test)
      next_stable_patch="v${s_major}.${s_minor}.$((s_patch + 1))"
      latest_test=$(git tag -l "${next_stable_patch}-test*" | grep -E "^${next_stable_patch}-test[0-9]+$" | sort -V | tail -n 1 || true)
      if [[ -n "$latest_test" ]]; then
          test_num=$(echo "$latest_test" | grep -oE "test[0-9]+" | sed 's/test//')
          next_tag="${next_stable_patch}-test$((test_num + 1))"
      else
          next_tag="${next_stable_patch}-test1"
      fi
      ;;
    *) echo "Unsupported release mode: $MODE" >&2; exit 1 ;;
  esac

  if [[ ! "$next_tag" =~ $VALID_TAG_REGEX ]]; then
    echo "Invalid tag calculated: $next_tag" >&2
    exit 1
  fi

  if git rev-parse "$next_tag" >/dev/null 2>&1; then
    REMOTE_TAG_SHA=$(git ls-remote --tags origin "refs/tags/$next_tag^{}" | awk '{print $1}')
    if [[ -z "$REMOTE_TAG_SHA" ]]; then
        REMOTE_TAG_SHA=$(git ls-remote --tags origin "refs/tags/$next_tag" | awk '{print $1}')
    fi
    if [[ "$REMOTE_TAG_SHA" == "${GITHUB_SHA}" ]]; then
      echo "Tag $next_tag exists and points to GITHUB_SHA. Safe to retry." >&2
      echo "$next_tag"
      exit 0
    else
      echo "Tag already exists and points to $REMOTE_TAG_SHA (expected ${GITHUB_SHA}): $next_tag" >&2
      echo "To retry a failed publication for this exact version, ensure release_version_override is used and GITHUB_SHA matches." >&2
      exit 1
    fi
  else
    git tag "$next_tag" "${GITHUB_SHA}"
    git push origin "$next_tag" >/dev/null 2>&1 || (
      VERIFY_SHA=$(git ls-remote --tags origin "refs/tags/$next_tag" | awk '{print $1}')
      if [[ "$VERIFY_SHA" == "${GITHUB_SHA}" ]]; then
        echo "Tag successfully verified on remote after push error." >&2
      else
        echo "Tag push failed and remote verification failed." >&2
        exit 1
      fi
    )
    echo "$next_tag"
    exit 0
  fi
fi
