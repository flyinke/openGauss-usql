# openGauss Fork Notes

This repository is a small fork of upstream `usql` for openGauss environments.

## Goals

- Keep the fork delta as small as possible.
- Preserve upstream `usql` structure so future syncs are easier.
- Produce portable Linux delivery artifacts for `amd64` and `arm64`.

## Versioning

Use the scheme:

```text
<upstream>-og.<n>
```

Example:

```text
0.21.4-og.1
```

Meaning:

- `0.21.4`: upstream `usql` base version
- `og`: openGauss fork marker
- `.1`: first fork release on that upstream base

Recommended follow-up examples:

- `0.21.4-og.2`: second fork release on the same upstream base
- `0.21.5-og.1`: first fork release after rebasing to upstream `0.21.5`

## Driver Notes

This fork replaces the PostgreSQL driver package with:

```text
gitee.com/opengauss/openGauss-connector-go-pq
```

The CLI supports the `opengauss://` scheme and also works with PostgreSQL-style URLs such as `postgres://`.

## Build Notes

For portable release artifacts in this fork, build only the PostgreSQL/openGauss path and disable optional chart rendering:

```sh
# x86_64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -trimpath \
  -ldflags='-s -w -X github.com/xo/usql/text.CommandName=gausssql -X github.com/xo/usql/text.CommandVersion=0.21.4-og.1' \
  -tags 'postgres no_base no_chart' \
  -o build/linux/amd64/0.21.4-og.1/gausssql .

# arm64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod=mod -trimpath \
  -ldflags='-s -w -X github.com/xo/usql/text.CommandName=gausssql -X github.com/xo/usql/text.CommandVersion=0.21.4-og.1' \
  -tags 'postgres no_base no_chart' \
  -o build/linux/arm64/0.21.4-og.1/gausssql .
```

Notes:

- Delivery binary name in this fork is `gausssql`.
- `no_chart` is currently used to avoid CGO-based chart rendering during cross-compilation.
- Output layout intentionally follows upstream `usql` conventions: `build/linux/<arch>/<version>/`.

## Connection Examples

```sh
# openGauss
gausssql opengauss://user:pass@host:port/postgres?sslmode=disable

# PostgreSQL-style URL also works in this fork
gausssql postgres://user:pass@host:port/postgres?sslmode=disable
```

## Release Guidance

Before publishing a fork release:

1. Rebuild both `amd64` and `arm64` artifacts.
2. Verify `--version` output matches the intended fork version.
3. Verify `opengauss://` is recognized and reaches the driver connection path.
4. Keep README changes minimal and prefer documenting fork-only behavior here.

## Automation

This fork keeps upstream workflows intact and adds fork-specific automation:

- `.github/workflows/opengauss-ci.yml`: builds `amd64` and `arm64` artifacts on push / pull request and runs a lightweight smoke test.
- `.github/workflows/release-opengauss.yml`: builds Linux release artifacts for `amd64` and `arm64` and creates a draft GitHub release.
- `build-opengauss.sh`: local and CI-friendly build entrypoint for release-style fork artifacts.

Recommended release flow:

1. Choose the next fork version, such as `0.21.4-og.2`.
2. Push a tag named `v0.21.4-og.2`.
3. Let `release-opengauss.yml` build both Linux artifacts and open a draft release.
4. Review attached files and publish the draft release.
