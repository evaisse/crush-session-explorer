from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import sqlite3
from typing import Iterable, Optional
import json


@dataclass
class Session:
    id: str
    title: Optional[str]
    created_at: Optional[str]
    metadata: Optional[str]
    content: Optional[str]
    message_count: Optional[int] = None


def list_messages(conn: sqlite3.Connection, session_id: str) -> list[dict]:
    cur = conn.execute(
        """
        SELECT id, role, parts, model, provider, created_at
        FROM messages
        WHERE session_id = ?
        ORDER BY created_at ASC
        """,
        (session_id,),
    )
    rows = cur.fetchall()
    msgs = []
    for r in rows:
        parts = []
        try:
            parts = json.loads(r["parts"]) if r["parts"] else []
        except Exception:
            parts = []
        msgs.append(
            {
                "id": r["id"],
                "role": r["role"],
                "parts": parts,
                "model": r["model"],
                "provider": r["provider"],
                "created_at": r["created_at"],
            }
        )
    return msgs


def connect(db_path: Path) -> sqlite3.Connection:
    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row
    return conn


def fetch_session(conn: sqlite3.Connection, session_id: str) -> Session:
    cur = conn.execute(
        """
        SELECT id, title, created_at, message_count
        FROM sessions
        WHERE id = ?
        """,
        (session_id,),
    )
    row = cur.fetchone()
    if row is None:
        raise KeyError(f"session not found: {session_id}")
    return Session(
        id=str(row["id"]),
        title=row["title"],
        created_at=str(row["created_at"]) if row["created_at"] is not None else None,
        metadata=None,
        content=None,
        message_count=int(row["message_count"]) if "message_count" in row.keys() and row["message_count"] is not None else None,
    )


def list_sessions(conn: sqlite3.Connection, limit: int = 20) -> list[Session]:
    cur = conn.execute(
        """
        SELECT id, title, created_at, message_count
        FROM sessions
        ORDER BY created_at DESC
        LIMIT ?
        """,
        (limit,),
    )
    rows = cur.fetchall()
    return [
        Session(
            id=str(r["id"]),
            title=r["title"],
            created_at=str(r["created_at"]) if r["created_at"] is not None else None,
            metadata=None,
            content=None,
            message_count=int(r["message_count"]) if "message_count" in r.keys() and r["message_count"] is not None else None,
        )
        for r in rows
    ]
