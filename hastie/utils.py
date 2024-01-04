from typing import List

# import operator
import re
import time


def urljoin(parts: List) -> str:
    """Join a list of URL paths avoiding duplicate /'s"""

    # remove leading and trailing slashes
    parts = list(map(lambda s: s.strip("/"), parts))

    # filter out empty
    parts = list(filter(lambda p: p != "", parts))

    # come together
    url = "/" + "/".join(parts) + "/"

    return url


def tryint(s):
    """Return an int if possible, or `s` unchanged."""
    try:
        return int(s)
    except ValueError:
        return s


def alphanum(s) -> List:
    """
    Turn a string into a list of string and number chunks.
    >>> alphanum_key("z23a")
    ["z", 23, "a"]
    """

    return [tryint(c) for c in re.split("([0-9]+)", s)]


def human_sort(k: List, field: str) -> List:
    k.sort(key=lambda el: alphanum(el[field]))
    return k


def date_sort(k: List) -> List:
    kd = list(filter(lambda el: "date" in el, k))
    kd.sort(key=lambda el: el["date"])
    kd.reverse()
    return kd


def timer(st: float = 0, s=""):
    if st == 0:
        return time.time()

    print(f"Elapsed for {s}: {time.time()-st:.1f}")
