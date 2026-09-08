"""Python launcher for the Go-powered yaju-bot."""

from __future__ import annotations

import hashlib
import os
import platform
import shutil
import stat
import subprocess
import sys
import tarfile
import tempfile
import urllib.error
import urllib.request
import zipfile
from importlib.metadata import PackageNotFoundError, version
from pathlib import Path
from typing import Optional

__all__ = ["run"]

_PACKAGE = "yaju-bot"
_REPO = "dtmpm3485/yaju-bot"
_FALLBACK_VERSION = "0.1.1"

try:
    __version__ = version(_PACKAGE)
except PackageNotFoundError:
    __version__ = _FALLBACK_VERSION


def _cache_root() -> Path:
    override = os.environ.get("YAJU_BOT_CACHE")
    if override:
        return Path(override).expanduser()
    if os.name == "nt" and os.environ.get("LOCALAPPDATA"):
        return Path(os.environ["LOCALAPPDATA"]) / "yaju-bot" / "cache"
    xdg = os.environ.get("XDG_CACHE_HOME")
    if xdg:
        return Path(xdg) / "yaju-bot"
    return Path.home() / ".cache" / "yaju-bot"


def _platform_asset_for(system: str, machine: str) -> tuple[str, str]:
    system = system.lower()
    machine = machine.lower()
    arch_map = {
        "x86_64": "amd64",
        "amd64": "amd64",
        "aarch64": "arm64",
        "arm64": "arm64",
        "armv7l": "armv7",
        "armv7": "armv7",
    }
    arch = arch_map.get(machine)
    if arch is None:
        raise RuntimeError(f"unsupported CPU architecture: {machine}")

    if system == "linux":
        os_name = "linux"
        ext = ".tar.gz"
    elif system == "android":
        os_name = "android"
        ext = ".tar.gz"
    elif system == "darwin":
        if arch == "armv7":
            raise RuntimeError("unsupported macOS architecture: armv7")
        os_name = "darwin"
        ext = ".tar.gz"
    elif system == "windows":
        if arch == "armv7":
            raise RuntimeError("unsupported Windows architecture: armv7")
        os_name = "windows"
        ext = ".zip"
    else:
        raise RuntimeError(f"unsupported operating system: {system}")

    return f"yaju-bot-{os_name}-{arch}{ext}", "yaju-bot.exe" if os_name == "windows" else "yaju-bot"


def _platform_asset() -> tuple[str, str]:
    return _platform_asset_for(platform.system(), platform.machine())


def _request_bytes(url: str) -> bytes:
    request = urllib.request.Request(url, headers={"User-Agent": f"yaju-bot/{__version__}"})
    with urllib.request.urlopen(request, timeout=30) as response:
        return response.read()


def _expected_checksum(checksums: bytes, asset_name: str) -> str:
    for raw_line in checksums.decode("utf-8").splitlines():
        parts = raw_line.strip().split()
        if len(parts) >= 2 and parts[-1].lstrip("*") == asset_name:
            digest = parts[0].lower()
            if len(digest) == 64 and all(ch in "0123456789abcdef" for ch in digest):
                return digest
    raise RuntimeError(f"checksum not found for release asset: {asset_name}")


def _download_release_binary() -> Path:
    asset_name, binary_name = _platform_asset()
    cache_dir = _cache_root() / __version__
    binary = cache_dir / binary_name
    if binary.is_file():
        return binary

    tag = f"v{__version__}"
    base = f"https://github.com/{_REPO}/releases/download/{tag}"
    try:
        checksums = _request_bytes(f"{base}/checksums.txt")
        archive = _request_bytes(f"{base}/{asset_name}")
    except (urllib.error.URLError, TimeoutError, OSError) as exc:
        raise RuntimeError(f"failed to download yaju-bot {tag}: {exc}") from exc

    expected = _expected_checksum(checksums, asset_name)
    actual = hashlib.sha256(archive).hexdigest()
    if actual != expected:
        raise RuntimeError(f"release checksum mismatch for {asset_name}")

    cache_dir.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory(prefix="yaju-bot-") as tmp_name:
        tmp = Path(tmp_name)
        archive_path = tmp / asset_name
        archive_path.write_bytes(archive)
        extracted = tmp / binary_name

        if asset_name.endswith(".zip"):
            with zipfile.ZipFile(archive_path) as zf:
                try:
                    data = zf.read(binary_name)
                except KeyError as exc:
                    raise RuntimeError(f"{binary_name} is missing from {asset_name}") from exc
                extracted.write_bytes(data)
        else:
            with tarfile.open(archive_path, "r:gz") as tf:
                try:
                    member = tf.getmember(binary_name)
                except KeyError as exc:
                    raise RuntimeError(f"{binary_name} is missing from {asset_name}") from exc
                source = tf.extractfile(member)
                if source is None:
                    raise RuntimeError(f"failed to read {binary_name} from {asset_name}")
                extracted.write_bytes(source.read())

        staged = cache_dir / f".{binary_name}.{os.getpid()}.tmp"
        try:
            shutil.copyfile(extracted, staged)
            if os.name != "nt":
                staged.chmod(staged.stat().st_mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)
            os.replace(staged, binary)
        finally:
            try:
                staged.unlink()
            except FileNotFoundError:
                pass

    return binary


def _source_checkout() -> Optional[Path]:
    here = Path(__file__).resolve()
    for parent in here.parents:
        if (parent / "go.mod").is_file() and (parent / "cmd" / "yaju-bot").is_dir():
            return parent
    return None


def _resolve_command() -> tuple[list[str], Optional[str]]:
    override = os.environ.get("YAJU_BOT_BIN")
    if override:
        return [override], None

    source = _source_checkout()
    go = shutil.which("go")
    if source is not None and go is not None:
        return [go, "run", "./cmd/yaju-bot"], str(source)

    try:
        return [str(_download_release_binary())], None
    except RuntimeError as download_error:
        if go is not None:
            return [go, "run", f"github.com/{_REPO}/cmd/yaju-bot@v{__version__}"], None
        raise RuntimeError(
            f"{download_error}. Install Go as a fallback, or set YAJU_BOT_BIN to a yaju-bot executable."
        ) from download_error


def run(token: str, *, data_dir: str | os.PathLike[str] | None = None) -> None:
    """Start yaju-bot and block until it exits."""
    if not isinstance(token, str) or not token.strip():
        raise ValueError("token must be a non-empty Discord bot token")

    env = os.environ.copy()
    env["DISCORD_TOKEN"] = token.strip()
    if data_dir is not None:
        env["YAJU_DATA_DIR"] = os.fspath(data_dir)

    command, cwd = _resolve_command()
    try:
        completed = subprocess.run(command, env=env, cwd=cwd, check=False)
    except KeyboardInterrupt:
        return
    if completed.returncode != 0:
        raise RuntimeError(f"yaju-bot exited with status {completed.returncode}")
