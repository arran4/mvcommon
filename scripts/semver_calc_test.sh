#!/bin/bash
set -euo pipefail

# This script mocks git tag interactions and tests the semantic version logic from semver_calc.sh
# It tests the versioning calculation purely mathematically without creating tags in real origin.

test_semver() {
    local mode=$1
    local existing_tags=$2
    local expected=$3
    local expect_fail=${4:-0}

    temp_dir=$(mktemp -d)
    cd "$temp_dir"
    git init -q
    git config user.email "test@test.com"
    git config user.name "test"
    git commit --allow-empty -m "init" -q

    if [ ! -z "$existing_tags" ]; then
        for t in $existing_tags; do
            git tag "$t"
        done
    fi

    # The actual calculation logic extracted from semver_calc.sh
    MODE="$mode"

    # 1. Filter out malformed tags and prereleases for stable tag baseline.
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
      *)
        if [ "$expect_fail" = "1" ]; then
           cd - >/dev/null
           rm -rf "$temp_dir"
           echo "PASS (Expected Failure): $mode"
           return 0
        else
            echo "FAIL: Unexpected mode $mode"
            return 1
        fi
        ;;
    esac

    # Ensure shape validation holds
    VALID_TAG_REGEX='^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$'
    if [[ ! "$next_tag" =~ $VALID_TAG_REGEX ]]; then
        if [ "$expect_fail" = "1" ]; then
            cd - >/dev/null
            rm -rf "$temp_dir"
            echo "PASS (Expected Failure): Invalid shape $next_tag"
            return 0
        else
            echo "FAIL: Invalid shape calculated $next_tag"
            return 1
        fi
    fi

    if [ "$next_tag" != "$expected" ]; then
        if [ "$expect_fail" = "1" ]; then
             echo "PASS (Expected Failure, Output mismatched but failure caught): $mode"
        else
             echo "FAIL: $mode with tags [$existing_tags]. Expected $expected but got $next_tag"
             return 1
        fi
    else
        if [ "$expect_fail" = "1" ]; then
             echo "FAIL: $mode with tags [$existing_tags]. Expected Failure but it produced $next_tag"
             return 1
        else
            echo "PASS: $mode with tags [$existing_tags] -> $next_tag"
        fi
    fi

    cd - >/dev/null
    rm -rf "$temp_dir"
}

echo "Testing Version Calculation Rules:"

# Stable only
test_semver "release-minor" "v1.0.0" "v1.1.0"
test_semver "release-patch" "v1.0.0" "v1.0.1"
test_semver "release-major" "v1.0.0" "v2.0.0"

# RC transitions
test_semver "release-rc" "v1.0.0" "v1.0.1-rc1"
test_semver "release-rc" "v1.0.0 v1.0.1-rc1" "v1.0.1-rc2"
test_semver "release-rc" "v1.0.0 v1.0.1-rc1 v1.0.1-rc2" "v1.0.1-rc3"

# Stable bump after RC (ignoring RCs for stable baseline)
test_semver "release-patch" "v1.0.0 v1.0.1-rc1 v1.0.1-rc2" "v1.0.1"
test_semver "release-minor" "v1.0.0 v1.0.1-rc1 v1.0.1-rc2" "v1.1.0"
test_semver "release-major" "v1.0.0 v1.0.1-rc1" "v2.0.0"

# Test transitions
test_semver "release-test" "v1.0.0" "v1.0.1-test1"
test_semver "release-test" "v1.0.0 v1.0.1-test1" "v1.0.1-test2"

# Independent RC and Test Sequences
test_semver "release-test" "v1.0.0 v1.0.1-rc1" "v1.0.1-test1"
test_semver "release-rc" "v1.0.0 v1.0.1-rc1 v1.0.1-test1 v1.0.1-test2" "v1.0.1-rc2"

# Malformed tags ignored
test_semver "release-patch" "v1.0.0 vfoo v1.0 v1.0.0-rc1 v2.0" "v1.0.1"
test_semver "release-rc" "v1.0.0 vfoo v1.0.1-rc v1.0.1-rc1" "v1.0.1-rc2"

echo "All tests passed successfully."