# crush-session-explorer

Small Python CLI to export Crush sessions from a SQLite DB to Markdown.

## Features
- Query sessions from sqlite (stdlib only)
- Render Markdown with YAML frontmatter
- Simple CLI: `python -m crush_md export`
- Tested with pytest, formatted/linted with ruff, typed with pyright

## Requirements
- Python 3.10+
- sqlite3 (stdlib)

Optional tools for dev:
- ruff, pyright, pytest

Install dev tools:
```
python -m venv .venv && . .venv/bin/activate && python -m pip install -U pip
pip install -r requirements.txt || true
pip install ruff pyright pytest
```

## Usage
Export a session to Markdown:
```
python -m crush_md export --db ./.crush/crush.db --session <id> --out out.md
```
Args:
- `--db`: Path to sqlite db (default .crush/crush.db)
- `--session`: Session id to export (required)
- `--out`: Output markdown file (required)

## Project layout
- `crush_md/db.py`: sqlite access and Session model
- `crush_md/markdown.py`: Markdown rendering
- `crush_md/cli.py`: argparse CLI entrypoint
- `tests/`: pytest tests

## Development
Run tests:
```
pytest -q
```
Lint and format:
```
ruff check .
ruff format .
```
Typecheck:
```
pyright
```

## Notes
- DB uses sqlite3.Row and parameterized queries only
- Markdown frontmatter lines <= 100 chars where possible
