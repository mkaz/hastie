from typing import Any

# import operator
import re
import time


def urljoin(parts: list[str]) -> str:
    """Join a list of URL paths avoiding duplicate /'s"""

    # remove leading and trailing slashes
    parts = list(map(lambda s: s.strip("/"), parts))

    # filter out empty
    parts = list(filter(lambda p: p != "", parts))

    # come together
    url = "/" + "/".join(parts) + "/"

    return url


def tryint(s: str) -> int | str:
    """Return an int if possible, or `s` unchanged."""
    try:
        return int(s)
    except ValueError:
        return s


def alphanum(s: str) -> list[int | str]:
    """
    Turn a string into a list of string and number chunks.
    >>> alphanum_key("z23a")
    ["z", 23, "a"]
    """

    return [tryint(c) for c in re.split("([0-9]+)", s)]


def human_sort(k: list[dict[str, Any]], field: str) -> list[dict[str, Any]]:
    k.sort(key=lambda el: alphanum(el[field]))
    return k


def date_sort(k: list[dict[str, Any]]) -> list[dict[str, Any]]:
    kd = list(filter(lambda el: "date" in el, k))
    kd.sort(key=lambda el: el["date"])
    kd.reverse()
    return kd


def timer(st: float = 0, s=""):
    if st == 0:
        return time.time()

    print(f"Elapsed for {s}: {time.time() - st:.1f}")
