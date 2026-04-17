# SortyMcSortFace

A tiny Go CLI that organizes files inside the `sort/` folder into type-based subfolders.

## What It Does

When you run the app, it scans files directly inside `sort/` and moves them into:

- `sort/images` for `.jpg`, `.jpeg`, `.png`, `.gif`
- `sort/videos` for `.mp4`, `.mkv`, `.avi`
- `sort/docs` for `.pdf`, `.docx`, `.txt`
- `sort/music` for `.mp3`, `.wav`
- `sort/others` for anything else

## Requirements

- Go `1.26.1` (or compatible with your local toolchain)

## Run

From the project root:

```bash
go run .
```

## Quick Example

1. Put some files into `sort/`, for example:

```text
sort/
	photo.jpg
	song.mp3
	notes.txt
	archive.zip
```

2. Run:

```bash
go run .
```

3. Result:

```text
sort/
	images/photo.jpg
	music/song.mp3
	docs/notes.txt
	others/archive.zip
```

## Notes

- Only top-level files in `sort/` are processed.
- Existing subfolders in `sort/` are skipped.
- Files are moved (renamed), not copied.
