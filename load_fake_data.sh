#!/bin/bash

go run . reset
go run . register rob
go run . addfeed "boot.dev Blog" "https://blog.boot.dev/index.xml"
go run . register nick
go run . login rob
go run . addfeed "TechCrunch" "https://techcrunch.com/feed/"
go run . addfeed "Hacker News" "https://news.ycombinator.com/rss"
go run . register matt
go run . follow "https://techcrunch.com/feed/"
go run . register phil
go run . follow "https://news.ycombinator.com/rss"
