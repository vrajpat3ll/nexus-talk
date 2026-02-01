#!/usr/bin/env python3
"""
Simple PDF to text extractor.

Usage:
  python tools/extract_pdfs.py path/to/file1.pdf [path/to/file2.pdf ...]

The script will try to use `pypdf` (preferred) or fall back to `pdfminer.six` if available.
If neither is installed it will print instructions to install them. This script does not create
or require a virtual environment; install packages globally if you choose.

Output files go into `extracted/` next to this script.
"""
from __future__ import annotations

import os
import sys
from pathlib import Path
from typing import List


def extract_with_pypdf(path: Path) -> str:
    from pypdf import PdfReader

    reader = PdfReader(str(path))
    texts: List[str] = []
    for page in reader.pages:
        try:
            texts.append(page.extract_text() or "")
        except Exception:
            # keep going if a page fails
            texts.append("")
    return "\n".join(texts)


def extract_with_pdfminer(path: Path) -> str:
    # pdfminer.six high-level API
    from pdfminer.high_level import extract_text

    return extract_text(str(path))


def main(argv: List[str]) -> int:
    if len(argv) < 2:
        print("Usage: python tools/extract_pdfs.py <file1.pdf> [file2.pdf ...]")
        return 2

    script_dir = Path(__file__).resolve().parent
    out_dir = script_dir.parent / "extracted"
    out_dir.mkdir(parents=True, exist_ok=True)

    # decide which backend is available
    backend = None
    try:
        import pypdf  # type: ignore

        backend = "pypdf"
    except Exception:
        try:
            import pdfminer  # type: ignore

            backend = "pdfminer"
        except Exception:
            backend = None

    if backend is None:
        print("No PDF extraction libraries found.")
        print("Install one of: pypdf (recommended) or pdfminer.six")
        print()
        print("Example (global install, no venv):")
        print("  pip install pypdf pdfminer.six")
        return 3

    for p in argv[1:]:
        pdf_path = Path(p)
        if not pdf_path.exists():
            print(f"Skipping: {pdf_path} (not found)")
            continue

        out_name = pdf_path.stem + ".txt"
        out_path = out_dir / out_name
        print(f"Extracting {pdf_path} -> {out_path} using {backend}...")
        try:
            if backend == "pypdf":
                text = extract_with_pypdf(pdf_path)
            else:
                text = extract_with_pdfminer(pdf_path)
        except Exception as e:
            print(f"Failed to extract {pdf_path}: {e}")
            continue

        # sanitize line endings
        text = text.replace('\r\n', '\n')

        with out_path.open("w", encoding="utf-8") as f:
            f.write(text)

        print(f"Wrote {out_path} ({len(text)} bytes)")

    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
