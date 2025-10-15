from __future__ import annotations

from datetime import datetime
from typing import Optional, List, Dict, Any
import json

from .db import Session


def _yaml_escape(s: str) -> str:
    return s.replace("\n", " ").replace('"', "'")


def _fmt_ts(ts: Any) -> str:
    if ts is None:
        return ""
    try:
        return datetime.fromtimestamp(int(ts)).strftime("%Y-%m-%d %H:%M")
    except Exception:
        try:
            return datetime.fromisoformat(str(ts).replace("Z", "+00:00")).astimezone().strftime("%Y-%m-%d %H:%M")
        except Exception:
            return str(ts)


def render_markdown(session: Session) -> str:
    title = session.title or f"Session {session.id}"
    created = session.created_at
    try:
        if created:
            created = datetime.fromtimestamp(int(created)).isoformat()
    except Exception:
        try:
            if created:
                created = datetime.fromisoformat(str(created).replace("Z", "+00:00")).isoformat()
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
    if session.message_count is not None:
        frontmatter_lines.append(f"message_count: {session.message_count}")
    if meta_obj is not None:
        frontmatter_lines.append("metadata:")
        for k, v in meta_obj.items():
            frontmatter_lines.append(f"  {k}: {json.dumps(v, ensure_ascii=False)}")
    frontmatter_lines.append("---\n")

    body = ""
    try:
        msgs: List[Dict[str, Any]] = json.loads(session.content) if session.content and session.content.strip().startswith("[") else None
    except Exception:
        msgs = None

    if isinstance(msgs, list):
        lines: List[str] = []
        for m in msgs:
            role = m.get("role", "")
            ts = _fmt_ts(m.get("created_at"))
            model = m.get("model") or ""
            provider = m.get("provider") or ""
            head = f"## {role} — {ts}"
            if model or provider:
                meta = "/".join([x for x in [model, provider] if x])
                head += f" ({meta})"
            lines.append(head)
            parts = m.get("parts") or []
            lines.append("<div>")
            for p in parts:
                if isinstance(p, str):
                    lines.append(p)
                elif isinstance(p, dict) and "text" in p:
                    lines.append(str(p["text"]))
            lines.append("</div>")
            lines.append("")
        body = "\n".join(lines)
    else:
        body = session.content or ""

    return "\n".join(frontmatter_lines) + body + ("\n" if body and not body.endswith("\n") else "")
