# SortyMcSortFace

A configurable Go CLI that organizes files into smart category folders.

It now uses multiple signals to classify files:

- Extension matching (large built-in list)
- Compound extension matching (`.tar.gz`, `.tar.bz2`, etc.)
- Filename keyword hints (for files like invoices, screenshots, backups)
- MIME/content sniffing fallback when extension is missing or unclear

## Features

- Organizes files into category folders under a target directory
- Handles duplicate filenames safely (`file_(1).ext`, `file_(2).ext`, ...)
- Optional recursive mode to process nested folders
- Dry-run mode to preview actions without moving anything
- Verbose plan output and end-of-run summary stats

## Categories

Built-in categories include:

- `images`
- `videos`
- `music`
- `docs`
- `presentations`
- `spreadsheets`
- `archives`
- `code`
- `installers`
- `fonts`
- Keyword-only buckets: `finance`, `screenshots`, `design`, `certificates`, `backups`
- Fallback: `others`

## Requirements

- Go `1.26.1` (or any compatible local Go toolchain)

## Usage

Basic run:

```bash
go run .
```

With flags:

```bash
go run . --dir sort --recursive --dry-run
```

## CLI Flags

- `--dir` (default: `sort`): directory to organize
- `--recursive` (default: `false`): process files in subfolders too
- `--dry-run` (default: `false`): show what would happen without moving files
- `--verbose` (default: `true`): print per-file action plan

## How Classification Works

For each file, SortyMcSortFace uses this order:

1. Compound extension (`.tar.gz`, `.tar.bz2`, `.tar.xz`, `.tar.zst`)
2. Normal extension lookup (e.g., `.png`, `.mp3`, `.xlsx`, `.go`)
3. Filename keyword hints (`invoice`, `screenshot`, `backup`, etc.)
4. Content-type sniffing (first 512 bytes via MIME detection)
5. Fallback to `others`

## Example

Input:

```text
sort/
  IMG_4102
  invoice_march.pdf
  project_backup_2026
  archive.tar.gz
  design_mockup_final
```

Command:

```bash
go run . --dir sort --dry-run
```

Possible result plan:

```text
Plan: sort/IMG_4102 -> sort/images/IMG_4102
Plan: sort/invoice_march.pdf -> sort/docs/invoice_march.pdf
Plan: sort/project_backup_2026 -> sort/backups/project_backup_2026
Plan: sort/archive.tar.gz -> sort/archives/archive.tar.gz
Plan: sort/design_mockup_final -> sort/design/design_mockup_final
```

## Creative Ideas To Build Next

1. Auto-watch mode
   Use `fsnotify` to watch the folder and sort files instantly as they arrive.

2. Rule profiles
   Add `profiles/work.json`, `profiles/media.json`, `profiles/dev.json` and switch with `--profile`.

3. AI-assisted categorization
   For ambiguous files, call an LLM with filename + extracted text to choose the best category.

4. Semantic duplicate detection
   Hash files to detect duplicates and move them to `duplicates/` instead of keeping copies.

5. Time-based organization
   Add optional destination patterns like `images/2026/04` or `docs/2026-Q2`.

6. OCR and document intelligence
   Extract text from scanned PDFs/images, then auto-route contracts/invoices/resumes to custom folders.

7. Undo log
   Write a JSON transaction log for every run and support `--undo last`.

8. Smart conflict naming
   Use rich conflict patterns like `file (from Downloads, 2026-04-17).ext`.

9. Dashboard mode
   Expose a tiny local web UI to show category counts, trends, and recent moves.

10. Cross-device sync helper
    After sorting, auto-mirror selected categories to cloud folders (Drive/Dropbox/S3).

## Notes

- The tool moves files (rename), it does not copy.
- In `--dry-run`, nothing on disk changes.
- If a file is already in its target folder, it is skipped.
