import hastie


def test_urljoin_base_empty():
    parts = ["/", "", "foo", "/bar"]
    url = hastie.utils.urljoin(parts)
    assert url == "/foo/bar"


def test_urljoin_base_one():
    parts = ["/bar", "", "foo", "/bar"]
    url = hastie.utils.urljoin(parts)
    assert url == "/bar/foo/bar"


def test_urljoin_base_inner():
    parts = ["/bar", "foo/bar", "/bar"]
    url = hastie.utils.urljoin(parts)
    assert url == "/bar/foo/bar/bar"
