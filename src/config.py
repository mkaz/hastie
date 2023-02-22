import argparse
from pathlib import Path
import sys
import toml
from typing import Dict


def init_args(version: str) -> Dict:
    parser = argparse.ArgumentParser(description="hastie")
    parser.add_argument("-q", "--quiet", action="store_true")
    parser.add_argument("-v", "--version", action="store_true")
    args = vars(parser.parse_args())

    # Convention over configuration.
    # No override, look for file in the current directory.
    conffile = "./hastie.toml"

    if args["version"]:
        print(f"hastie v{version}")
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