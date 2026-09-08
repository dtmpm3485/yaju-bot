"""Python launcher for the Go-powered yaju-bot."""

from __future__ import annotations

import os
import shutil
import subprocess
from pathlib import Path
from typing import Optional

__all__ = ["run"]
__version__ = "0.1.0"


def _bundled_binary() -> Optional[Path]:
    name = "yaju-bot.exe" if os.name == "nt" else "yaju-bot"
    candidate = Path(__file__).with_name("bin") / name
    return candidate if candidate.is_file() else None


def run(token: str, *, data_dir: str | os.PathLike[str] | None = None) -> None:
    """Start yaju-bot.

    The Discord token is passed only through the child process environment.
    A bundled native binary is preferred. In a source checkout, Go is used as
    a fallback so development and Termux usage stay simple.
    """
    if not isinstance(token, str) or not token.strip():
        raise ValueError("token must be a non-empty Discord bot token")

    env = os.environ.copy()
    env["DISCORD_TOKEN"] = token.strip()
    if data_dir is not None:
        env["YAJU_DATA_DIR"] = os.fspath(data_dir)

    override = env.get("YAJU_BOT_BIN")
    if override:
        command = [override]
        cwd = None
    else:
        bundled = _bundled_binary()
        if bundled is not None:
            command = [str(bundled)]
            cwd = None
        else:
            go = shutil.which("go")
            if go is None:
                raise RuntimeError(
                    "yaju-bot native binary was not found and Go is not installed. "
                    "Install Go or set YAJU_BOT_BIN to a yaju-bot executable."
                )

            repo_root = Path(__file__).resolve().parents[2]
            if (repo_root / "go.mod").is_file():
                command = [go, "run", "./cmd/yaju-bot"]
                cwd = str(repo_root)
            else:
                command = [
                    go,
                    "run",
                    "github.com/dtmpm3485/yaju-bot/cmd/yaju-bot@latest",
                ]
                cwd = None

    try:
        completed = subprocess.run(command, env=env, cwd=cwd, check=False)
    except KeyboardInterrupt:
        return
    if completed.returncode != 0:
        raise RuntimeError(f"yaju-bot exited with status {completed.returncode}")
