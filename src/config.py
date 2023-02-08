import argparse
import os
from pathlib import Path
import sys
import toml
from typing import Dict

VERSION = "1.0.0"


def init_args() -> Dict:
    parser = argparse.ArgumentParser(description="hastie")
    parser.add_argument("-v", "--version", action="store_true")
    args = vars(parser.parse_args())

    # Allow override? Why complicate matters, just look for it
    # in the current directory. Always the same name.
    conffile = "./hastie.conf"

    if args["version"]:
        print(f"hastie v{VERSION}")
        sys.exit()

    # args is dict return, so add defaults
    args["content_dir"] = "./content"
    args["templates_dir"] = "./templates"
    args["output_dir"] = "./output"

    ## read config
    config = toml.load(conffile)

    # Merge config over args
    # anything in config will overwrite defaults
    args = args | config

    # Check content and template directories exists
    args["content_dir"] = Path(args["content_dir"])
    args["templates_dir"] = Path(args["templates_dir"])
    if not args["content_dir"].is_dir():
        print("Content directory '{}' not found".format(args["content_dir"]))
        sys.exit()

    if not args["templates_dir"].is_dir():
        print("Templates directory '{}' not found".format(args["templates_dir"]))
        sys.exit()

    return args
