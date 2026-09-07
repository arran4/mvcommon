#!/bin/bash
set -euo pipefail

# This script acts as a thin wrapper for repository-specific policy and transactional safety.
# It delegates all actual version arithmetic and parsing to the authoritative git-tag-inc CLI.

MODE="${RELEASE_MODE:-}"
OVERRIDE="${RELEASE_VERSION_OVERRIDE:-}"
GH_SHA="${GITHUB_SHA:-}"
MAIN_SHA="${ORIGIN_MAIN_SHA:-}"

if [[ "${GH_SHA}" != "${MAIN_SHA}" ]]; then
  echo "Error: Validation failed. ${GH_SHA} is not the current origin/main (${MAIN_SHA})" >&2
  exit 1
fi

# Ensure pinned git-tag-inc is installed so CI and production always use the identical version
go install github.com/arran4/git-tag-inc/cmd/git-tag-inc@90266586fefee6ffcb9fb02b00543b5959cd6c13
export PATH="$HOME/go/bin:$PATH"

VALID_TAG_REGEX='^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$'

if [[ -n "$OVERRIDE" ]]; then
  # Shape Validation
  OVERRIDE="${OVERRIDE#v}"
  next_tag="v$OVERRIDE"

  if [[ ! "$next_tag" =~ $VALID_TAG_REGEX ]]; then
      echo "Error: Override tag does not meet required shape: $next_tag" >&2
      exit 1
  fi

  # Existing Tag Validation
  if git rev-parse "$next_tag" >/dev/null 2>&1; then
    # Annotated dereferencing
    REMOTE_TAG_SHA=$(git ls-remote --tags origin "refs/tags/$next_tag^{}" | awk '{print $1}')
    if [[ -z "$REMOTE_TAG_SHA" ]]; then
        REMOTE_TAG_SHA=$(git ls-remote --tags origin "refs/tags/$next_tag" | awk '{print $1}')
    fi

    if [[ "$REMOTE_TAG_SHA" == "${GH_SHA}" ]]; then
      echo "Recovery Tag $next_tag exists and points to GITHUB_SHA. Safe to retry." >&2
    else
      echo "Recovery Tag already exists and points to $REMOTE_TAG_SHA (expected ${GH_SHA}): $next_tag" >&2
      exit 1
    fi
  else
    echo "Error: Recovery tag $next_tag does not exist. The override is for recovery only." >&2
    exit 1
  fi

else
  # Filter strictly for valid stable tags
  latest_stable=$(git tag -l "v*" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -n 1 || true)
  if [[ -z "$latest_stable" ]]; then latest_stable="v0.0.0"; fi

  case "$MODE" in
    release-major)
       next_tag=$(git-tag-inc -print-version-only -base-version "$latest_stable" major)
       ;;
    release-minor)
       next_tag=$(git-tag-inc -print-version-only -base-version "$latest_stable" minor)
       ;;
    release-patch)
       next_tag=$(git-tag-inc -print-version-only -base-version "$latest_stable" patch)
       ;;
    release-rc)
       next_patch=$(git-tag-inc -print-version-only -base-version "$latest_stable" patch)
       latest_rc=$(git tag -l "${next_patch}-rc*" | grep -E "^${next_patch}-rc[0-9]+$" | sort -V | tail -n 1 || true)
       if [[ -n "$latest_rc" ]]; then
           next_tag=$(git-tag-inc -print-version-only -base-version "$latest_rc" rc)
       else
           next_tag=$(git-tag-inc -print-version-only -base-version "$latest_stable" patch rc)
       fi
       ;;
    release-test)
       next_patch=$(git-tag-inc -print-version-only -base-version "$latest_stable" patch)
       latest_test=$(git tag -l "${next_patch}-test*" | grep -E "^${next_patch}-test[0-9]+$" | sort -V | tail -n 1 || true)
       if [[ -n "$latest_test" ]]; then
           next_tag=$(git-tag-inc -print-version-only -base-version "$latest_test" test)
       else
           next_tag=$(git-tag-inc -print-version-only -base-version "$latest_stable" patch test)
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
    if [[ "$REMOTE_TAG_SHA" == "${GH_SHA}" ]]; then
      echo "Tag $next_tag exists and points to GITHUB_SHA. Safe to retry." >&2
    else
      echo "Tag already exists and points to $REMOTE_TAG_SHA (expected ${GH_SHA}): $next_tag" >&2
      echo "To retry a failed publication for this exact version, ensure release_version_override is used and GITHUB_SHA matches." >&2
      exit 1
    fi
  else
    # Automatically apply the calculated tag
    git tag "$next_tag" "${GH_SHA}"
    git push origin "$next_tag" >/dev/null 2>&1 || (
      VERIFY_SHA=$(git ls-remote --tags origin "refs/tags/$next_tag" | awk '{print $1}')
      if [[ "$VERIFY_SHA" == "${GH_SHA}" ]]; then
        echo "Tag successfully verified on remote after push error." >&2
      else
        echo "Tag push failed and remote verification failed." >&2
        exit 1
      fi
    )
  fi
fi

# Print final result for workflow consumption
echo "$next_tag"
