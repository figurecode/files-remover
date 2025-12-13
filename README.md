# files-remover
![Go](https://github.com/figurecode/files-remover/actions/workflows/go.yml/badge.svg)
![GitHub release](https://img.shields.io/github/release/figurecode/files-remover.svg)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

[![English](https://img.shields.io/badge/Language-English-blue.svg)](README.md)
[![Русский](https://img.shields.io/badge/Language-Русский-blue.svg)](README.ru.md)

A Go utility for bulk deletion of files by name in a specified directory.

Useful for cleaning up temporary files, caches, and backups like:
- `backup-2025-01-01.zip`
- `session_12345.tmp`
- `debug-2025.log`

## Features

- Search by exact name or by prefix before a separator
- Exclude arbitrary subdirectories (e.g., `node_modules`, `.git`)
- Demo mode (`-m true`) — shows what will be deleted without touching anything
- Human-readable report: number of files and freed space (in KB/MB/GB)
- Safe deletion: "file not exists" errors are ignored (TOCTOU protection)

## Installation

The simplest way (requires Go 1.25+):

```bash
go install github.com/figurecode/files-remover/cmd/files-remover@latest
```

Or build from source:

```bash
git clone https://github.com/figurecode/files-remover.git
cd files-remover
go build -o files-remover cmd/files-remover/main.go
```

## Usage

### Main flags

```bash
./files-remover -d <directory> [flags] <pattern1> [pattern2...]
```

| Flag | Required?    | Description                                                                              | Default            |
|:-----|:------------|------------------------------------------------------------------------------------------|--------------------|
| `-d` | Yes*        | Path to the search directory (*required if not run from the target folder)               | Current directory  |
| `-e` | No          | Excluded subdirectories (comma-separated)                                                | (none)             |
| `-m` | No          | `true` — demo mode, `false` — real deletion                                              | `true`             |
| `-s` | No          | Separator for splitting filename parts                                                   | (none)             |

### Examples

1. Show what will be deleted (demo mode by default):

```bash
./files-remover -d /var/log app-debug.log backup-old.zip
```

2. Actually delete all files whose names start with `temp-2025` (e.g., temp-2025-12-12.log, temp-2025-11-01.log):

```bash
./files-remover -d /tmp -m false -s "-" temp-2025
```

3. Delete all `access-2024.log` and `error-2024.log` files, excluding folders `/var/log/journal`, `.snapshots`:

```bash
./files-remover -d /var/log -e /var/log/journal,.snapshots -m false access-2024.log error-2024.log
```

4. Delete everything starting with `cache-` in the name:

```bash
./files-remover -d ~/Library/Caches -m false -s "-" cache
```

To use a different separator (e.g., `_`):

```bash
./files-remover -d /data -s _ -m false session
```

## Demo mode output (example)

```text
127 files will be deleted in total
3.4 GB will be freed up

Files to be deleted:
---------------------------------
PATH: /tmp/backup-full-2025-01-01.tar.gz
---------------------------------
PATH: /tmp/session_8f3d9a21.tmp
---------------------------------
PATH: /var/log/nginx/access-2024-12-01.log

END
```

## Safety

- Demo mode by default
- Always run first without `-m false`
- "File already deleted" errors are ignored — the utility won't crash due to race conditions

## License

MIT

## Contributing

If you want to contribute to the project, create an issue or send a pull request on GitHub.

## Acknowledgments

Thanks to the following people for their contributions:

- [@etilite](https://github.com/etilite) — code review and valuable advice
