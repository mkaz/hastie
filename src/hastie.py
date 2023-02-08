#!/usr/bin/env python3

from config import init_args
from typing import List


def main():
    args = init_args()
    print("Welcome to Hastie")

    pages = args["content_dir"].glob("**/*.md")

    for page in pages:
        print(page)


if __name__ == "__main__":
    main()
