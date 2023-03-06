"""RSS Module"""
import operator
from datetime import datetime
from email.utils import formatdate
from textwrap import dedent
from typing import Dict, List


def generate_rss(config: Dict, pages: List) -> str:
    """With the config and list of pages, returns an RSS document"""

    # we want at least a slash for the link
    if "baseurl" in config:
        baseurl = config["baseurl"]
    else:
        baseurl = "/"

    rss = f"""<?xml version="1.0"?>
        <rss version="2.0">
            <channel>
                <title>{config["site"]["title"]}</title>
                <link>{baseurl}</link>
                <description>{config["site"]["description"]}</description>
                <language>en-us</language>
                <pubDate>{formatdate()}</pubDate>
                <lastBuildDate>{formatdate()}</lastBuildDate>
                <generator>Hastie</generator>
                <managingEditor>{config["site"]["author"]}</managingEditor>
                <webMaster>{config["site"]["author"]}</webMaster>
    """

    # filter out pages without a date
    pages = list(filter(lambda p: "date" in p, pages))

    # filter out draft pages
    pages = list(filter(lambda p: "draft" not in p, pages))

    # sort by date
    pages.sort(key=operator.itemgetter("date"), reverse=True)

    # limit to 10 most recent
    for page in pages[:10]:
        pubdate = datetime.combine(page["date"], datetime.min.time())
        rss += f"""
            <item>
                <title>{page["title"]}</title>
                <link>{page["url"]}</link>
                <description><![CDATA[
                    {page["content"]}
                ]]></description>
                <guid>{page["url"]}</guid>
                <pubDate>{formatdate(pubdate.timestamp())}</pubDate>
            </item>"""

    rss += "</channel></rss>"
    return dedent(rss)
