# Description

This PR implements a safe manual release architecture, solving the missing safe manual release path while preserving existing prerelease-safety fixes.

## Changes:
- **Manual Release Lifecycle:** Added `workflow_dispatch` trigger with modes (`release-major`, `release-minor`, `release-patch`, `release-rc`, `release-test`) and `publish-tag` mode for internal dispatch.
- **Validation Before Tagging:** A new `prepare-release` job runs `go test ./...`, GoReleaser checks, and snapshot builds. Validation happens *before* a permanent tag is created.
- **Exact-Commit Tagging:** The workflow calculates the semantic version, verifies remote state, and tags exactly the validated `${GITHUB_SHA}`.
- **Explicit Publisher Dispatch:** After pushing the tag, the workflow mechanically uses `gh workflow run` to invoke the `publish-tag` mode, preventing recursion issues associated with `GITHUB_TOKEN` tag pushes.
- **Non-recursion:** `publish-tag` requires an eligible tag context and does not create tags or recursively dispatch itself.
- **Recovery Behavior:** A `release_version_override` input is provided. If the intended tag already exists, the workflow verifies it resolves to the exact validated commit before safely continuing. If the SHA differs, the job fails safely.
- **External Tags:** External user tag pushes remain unchanged and continue to trigger the GoReleaser publication.
- **GoReleaser as Sole Owner:** GoReleaser remains the only creator/owner of GitHub Releases for any specific tag, satisfying single-owner invariant.
- **Permissions:** Permissions are constrained. The global workflow runs with `contents: read`, while jobs individually elevate `contents: write`, `packages: write`, and `actions: write` only when required.

## Validation Performed:
- Verified `go test ./...` passes.
- Verified GoReleaser config validity.
- Stable Docker releases include `latest`, while prereleases omit it (inherited from existing behavior).
- Pre-release Homebrew publication remains skipped (inherited).

References:
- https://arran4.github.io/blog/post/2026/042-simplified-github-ci-release-safe/
- https://arran4.github.io/blog/post/2026/041-release-safe-single-owner-github-ci/
