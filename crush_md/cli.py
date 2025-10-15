from __future__ import annotations

import argparse
from pathlib import Path
import sys
from datetime import datetime
import re

from .db import connect, fetch_session, list_sessions, list_messages
from .markdown import render_markdown


def slugify(text: str) -> str:
    text = text.strip().lower()
    text = re.sub(r"[^a-z0-9\-\s_]+", "", text)
    text = re.sub(r"[\s_]+", "-", text)
    return text or "untitled"


def local_dt(ts: str | None) -> str:
    if not ts:
        return ""
    try:
        dt = datetime.fromtimestamp(int(ts))
        return dt.strftime("%Y-%m-%d %H:%M")
    except Exception:
        try:
            return datetime.fromisoformat(ts.replace("Z", "+00:00")).astimezone().strftime("%Y-%m-%d %H:%M")
        except Exception:
            return ts


def main(argv: list[str] | None = None) -> int:
    p = argparse.ArgumentParser(prog="crush-md")
    sub = p.add_subparsers(dest="cmd", required=True)

    exp = sub.add_parser("export", help="Export session to markdown")
    exp.add_argument("--db", type=Path, default=Path(".crush/crush.db"))
    exp.add_argument("--session")
    exp.add_argument("--out", type=Path)

    args = p.parse_args(argv)

    if args.cmd == "export":
        conn = connect(args.db)
        session_id = args.session
        if not session_id:
            sessions = list_sessions(conn, 50)
            for i, s in enumerate(sessions, 1):
                print(f"{i:2d}. {s.id} — {local_dt(s.created_at)} — {s.title or ''} — {(s.message_count or 0)} msg")
            sel = input("Select session number: ").strip()
            if not sel.isdigit() or int(sel) < 1 or int(sel) > len(sessions):
                print("invalid selection", file=sys.stderr)
                return 2
            session_id = sessions[int(sel) - 1].id
        session = fetch_session(conn, session_id)
        msgs = list_messages(conn, session.id)
        parts_text = []
        for m in msgs:
            if m["parts"]:
                for p in m["parts"]:
                    if isinstance(p, str):
                        parts_text.append(p)
                    elif isinstance(p, dict) and "text" in p:
                        parts_text.append(str(p["text"]))
        session.content = "\n\n".join([t for t in parts_text if t])
        out_path = args.out
        if not out_path:
            base = slugify(session.title or f"session-{session.id[:8]}")
            prefix = datetime.fromtimestamp(int(session.created_at)).strftime("%Y-%m-%d_%H-%M") if session.created_at and str(session.created_at).isdigit() else datetime.now().strftime("%Y-%m-%d_%H-%M")
            default_dir = Path(".crush/sessions")
            default_name = f"{prefix}_{base}.md"
            default_path = default_dir / default_name
            entered = input(f"Output file path [{default_path}]: ").strip()
            out_path = Path(entered) if entered else default_path
        md = render_markdown(session)
        out_path.parent.mkdir(parents=True, exist_ok=True)
        out_path.write_text(md, encoding="utf-8")
        print(str(out_path))
        return 0

    return 1


if __name__ == "__main__":
    raise SystemExit(main())
