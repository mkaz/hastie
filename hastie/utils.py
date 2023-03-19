from typing import List


def urljoin(parts: List) -> str:
    """Join a list of URL paths avoiding duplicate /'s"""

    # remove leading and trailing slashes
    parts = list(map(lambda s: s.strip("/"), parts))

    # filter out empty
    parts = list(filter(lambda p: p != "", parts))

    # come together
    url = "/" + "/".join(parts) + "/"

    return url
