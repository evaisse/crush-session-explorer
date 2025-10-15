from __future__ import annotations

from datetime import datetime
from typing import Optional
import json

from .db import Session


def _yaml_escape(s: str) -> str:
    return s.replace("\n", " ").replace('"', "'")


def render_markdown(session: Session) -> str:
    title = session.title or f"Session {session.id}"
    created = session.created_at
    try:
        if created:
            # Normalize to ISO 8601
            created = datetime.fromisoformat(created.replace("Z", "+00:00")).isoformat()
    except Exception:
        pass
    meta_obj: Optional[dict] = None
    if session.metadata:
        try:
            meta_obj = json.loads(session.metadata)
        except Exception:
            meta_obj = {"raw": session.metadata}

    frontmatter_lines = ["---"]
    frontmatter_lines.append(f'title: "{_yaml_escape(title)}"')
    frontmatter_lines.append(f"session_id: {session.id}")
    if created:
        frontmatter_lines.append(f"created_at: {created}")
    if meta_obj is not None:
        frontmatter_lines.append("metadata:")
        for k, v in meta_obj.items():
            frontmatter_lines.append(f"  {k}: {json.dumps(v, ensure_ascii=False)}")
    frontmatter_lines.append("---\n")

    body = session.content or ""

    return "\n".join(frontmatter_lines) + body + ("\n" if not body.endswith("\n") else "")
