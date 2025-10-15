from __future__ import annotations

import sqlite3
from pathlib import Path

from crush_md.db import connect, fetch_session, Session


def test_fetch_session(tmp_path: Path):
    db = tmp_path / "test.db"
    conn = sqlite3.connect(db)
    conn.execute(
        "CREATE TABLE sessions (id TEXT PRIMARY KEY, title TEXT, created_at TEXT, metadata TEXT, content TEXT)"
    )
    conn.execute(
        "INSERT INTO sessions(id, title, created_at, metadata, content) VALUES (?, ?, ?, ?, ?)",
        ("id1", "Titre", "2025-10-01T12:00:00Z", "{\"k\":\"v\"}", "# body\n"),
    )
    conn.commit()
    conn.close()

    c = connect(db)
    s = fetch_session(c, "id1")

    assert isinstance(s, Session)
    assert s.id == "id1"
    assert s.title == "Titre"
