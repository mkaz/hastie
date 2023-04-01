import os
from pathlib import Path
import shutil
import sysrsync


def get_output_file(f: Path, c: Path, o: Path) -> Path:
    """Takes filename, content directory, output directory and returns output file"""
    # get file relative to content directory
    jf = Path(os.path.relpath(f, start=c))

    if jf.name == "index.md":
        return Path(o, jf.parent, "index.html")

    # filename without extension is .stem
    # create directory from stem - write as index.html
    return Path(o, jf.parent, jf.stem, "index.html")


def copy_static_assets(cdir: Path, odir: Path, tdir: Path):
    # copy templates static dir to output
    # - content at top level: for example /favicon.ico
    tpl_static = Path(tdir, "static")
    out_tpl_static = Path(odir)
    shutil.copytree(tpl_static, out_tpl_static, dirs_exist_ok=True)

    # copy site static dir to output
    # - content within under /static dir
    site_static = Path("./", "static")
    out_static = Path(odir, "static")
    if site_static.is_dir():
        shutil.copytree(site_static, out_static, dirs_exist_ok=True)

    # sync content structure to output excluding markdown
    sysrsync.run(
        source=str(cdir),  # add trailing slash
        destination=str(odir),  # add trailing slash
        options=["-a"],
        exclusions=["*.md"],
    )
