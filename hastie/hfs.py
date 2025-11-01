from pathlib import Path
import shutil


def get_output_file(f: Path, c: Path, o: Path) -> Path:
    """Takes filename, content directory, output directory and returns output file"""
    # get file relative to content directory
    jf = f.relative_to(c)

    if jf.name == "index.md":
        return Path(o, jf.parent, "index.html")

    # filename without extension is .stem
    # create directory from stem - write as index.html
    return Path(o, jf.parent, jf.stem, "index.html")


def copy_static_assets(cdir: Path, odir: Path, static_dir: Path):
    # copy configured static dir to output
    # - content at top level: for example /favicon.ico
    if static_dir.is_dir():
        shutil.copytree(static_dir, odir, dirs_exist_ok=True)

    # copy site static dir to output
    # - content within under /static dir
    site_static = Path("./", "static")
    out_static = Path(odir, "static")
    if site_static.is_dir():
        shutil.copytree(site_static, out_static, dirs_exist_ok=True)

    # sync content structure to output excluding markdown
    shutil.copytree(
        cdir,
        odir,
        ignore=shutil.ignore_patterns("*.md"),
        dirs_exist_ok=True,
    )
