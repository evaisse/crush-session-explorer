from __future__ import annotations

import argparse
from pathlib import Path
import sys

from .db import connect, fetch_session, list_sessions, list_messages
from .markdown import render_markdown


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
        if not args.session:
            sessions = list_sessions(conn, 20)
            for s in sessions:
                title = s.title or ""
                created = s.created_at or ""
                print(f"{s.id}\t{created}\t{title}\t{(s.message_count or 0)}")
            return 0
        if not args.out:
            print("--out is required when --session is provided", file=sys.stderr)
            return 2
        session = fetch_session(conn, args.session)
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
        md = render_markdown(session)
        args.out.write_text(md, encoding="utf-8")
        return 0

    return 1


if __name__ == "__main__":
    raise SystemExit(main())
