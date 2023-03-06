from typing import List


def urljoin(parts: List) -> str:
    if parts[0].startswith("/"):
        start = "/"
    else:
        start = ""

    # remove slash
    parts = list(map(slash_strip, parts))

    # filter out empty
    parts = list(filter(lambda p: p != "", parts))

    # add start to front of list
    parts.insert(0, start)

    # come together
    url = "/".join(parts)

    return url


def slash_strip(s: str) -> str:
    return s.strip("/")
