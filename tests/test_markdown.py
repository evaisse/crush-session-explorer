from __future__ import annotations

from crush_md.db import Session
from crush_md.markdown import render_markdown


def test_render_markdown_frontmatter_and_body():
    s = Session(
        id="abc123",
        title="Ma session",
        created_at="2025-10-01T12:34:56Z",
        metadata='{"foo":"bar"}',
        content="# Titre\nContenu\n",
    )
    md = render_markdown(s)
    assert md.startswith("---\n")
    assert "title: \"Ma session\"" in md
    assert "session_id: abc123" in md
    assert "created_at:" in md
    assert "metadata:" in md
    assert "foo: \"bar\"" in md
    assert md.endswith("\n")
