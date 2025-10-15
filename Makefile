PY=python3
PIP=$(PY) -m pip
VENV=.venv
ACT=. $(VENV)/bin/activate

.DEFAULT_GOAL := help

.PHONY: help venv install lint fmt typecheck test export build clean

help:
	@echo "Available targets:"
	@echo "  venv        Create virtualenv (.venv)"
	@echo "  install     Install dev deps (ruff, pyright, pytest, pyinstaller)"
	@echo "  lint        Run ruff lint"
	@echo "  fmt         Run ruff format"
	@echo "  typecheck   Run pyright"
	@echo "  test        Run pytest"
	@echo "  export      Export markdown (SESSION, OUT)"
	@echo "  build       Build standalone binary with PyInstaller -> dist/crush-md"
	@echo "  clean       Clean venv, caches, build artifacts"

venv:
	$(PY) -m venv $(VENV)
	$(ACT) && $(PY) -m pip install -U pip

install:
	$(ACT) && $(PIP) install -r requirements.txt || true
	$(ACT) && $(PIP) install ruff pyright pytest pyinstaller

lint:
	$(ACT) && ruff check .

fmt:
	$(ACT) && ruff format .

typecheck:
	$(ACT) && pyright

test:
	$(ACT) && pytest -q

export:
	$(ACT) && CMD="$(PY) -m crush_md.cli export --db ./.crush/crush.db"; \
	[ -n "$(SESSION)" ] && CMD="$$CMD --session $(SESSION)"; \
	[ -n "$(OUT)" ] && CMD="$$CMD --out $(OUT)"; \
	eval "$$CMD"

build:
	$(ACT) && pyinstaller -F -n crush-md -p . crush_md/cli.py
	@echo "Binary: dist/crush-md"

clean:
	rm -rf $(VENV) __pycache__ .pytest_cache .ruff_cache .pyright build dist *.spec
