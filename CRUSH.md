CRUSH.md

Stack
- Python 3.10+
- sqlite3 (stdlib), pathlib, argparse, dataclasses, typing
- Testing: pytest

Commands
- Setup venv: python -m venv .venv && . .venv/bin/activate && python -m pip install -U pip
- Install deps: pip install -r requirements.txt || true (stdlib only)
- Lint: ruff check .
- Format: ruff format .
- Typecheck: pyright
- Test: pytest -q
- Test single: pytest -q -k "<pattern>"

Conventions
- Imports: stdlib first, third-party next, local last; absolute imports preferred.
- Formatting: Ruff (PEP8-ish) with 100 char line length; use format command above.
- Types: annotate public functions; enable from __future__ import annotations.
- Naming: snake_case for vars/functions, PascalCase for classes, UPPER_SNAKE for consts.
- Errors: raise explicit exceptions; never swallow; log with logging not prints.
- I/O: use pathlib.Path and context managers; prefer UTF-8.
- DB: use sqlite3 with row_factory=sqlite3.Row; parameterized queries only.
- Markdown: keep lines <= 100, use fenced code blocks, YAML frontmatter when needed.

Cursor/Copilot
- No Cursor or Copilot rules found; if added later, mirror key rules here.

Repo Tasks
- CLI: python -m crush_md export --db ./.crush/crush.db --session <id> --out out.md
- Library: crush_md/db.py (query), crush_md/markdown.py (render), crush_md/cli.py (argparse)
- Tests: tests/test_markdown.py, tests/test_db.py
