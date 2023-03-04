from typing import List


def urljoin(parts: List) -> str:
    if parts[0].startswith("/"):
        start = "/"
    else:
        start = ""

    # remove slash
    parts = list(map(slash_strip, parts))

    # filter out empty
    parts = filter(lambda p: p != "", parts)

    # come together
    url = start + "/".join(parts)

    return url


def slash_strip(s: str) -> str:
    return s.strip("/")
