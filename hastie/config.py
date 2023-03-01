import argparse
from pathlib import Path
import sys
import toml

# initialize config
VERSION = "0.9.2"


parser = argparse.ArgumentParser(description="hastie")
parser.add_argument("-q", "--quiet", action="store_true")
parser.add_argument("-v", "--version", action="store_true")
parser.add_argument("-c", "--conf", help="Config file")
args = vars(parser.parse_args())

# Convention over configuration.
# No override, look for file in the current directory.
if args["conf"]:
    conffile = Path(args["conf"])
else:
    conffile = Path("./hastie.toml")

if not conffile.is_file():
    print("Error: hastie.toml file not found")
    print("Are you in the right directory?")
    sys.exit()

if args["version"]:
    print(f"hastie v{VERSION}")
    sys.exit()

# args is dict return, so add defaults
args["content_dir"] = "./content"
args["templates_dir"] = "./templates"
args["output_dir"] = "./output"

## read config
conf = toml.load(conffile)

# Merge config over args
# anything in config will overwrite defaults
config = args | conf

# Convert directories to Paths
config["content_dir"] = Path(config["content_dir"])
config["templates_dir"] = Path(config["templates_dir"])
