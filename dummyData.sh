#!/bin/zsh

go run . reset

go run . register Ant
go run . register Bird
go run . register Cat
go run . register Dog

go run . login Ant
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "LaneBlog" https://www.wagslane.dev/index.xml

go run . login Bird
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "MassivelyOP Site" https://massivelyop.com/feed/

go run . login Cat
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "MassivelyOP Site" https://massivelyop.com/feed/
go run . addFeed "Guild Wars 2" https://www.guildwars2.com/en/feed/

go run . login Ant
go run . follow https://feeds.libsyn.com/63150/rss

go run . login Bird
go run . follow https://www.guildwars2.com/en/feed/
go run . follow https://massivelyop.com/feed/
go run . follow https://www.wagslane.dev/index.xml

echo "Some sample feeds for manual testing:"
echo "HN | https://hnrss.org/newest"
echo "LaneBlog | https://www.wagslane.dev/index.xml"
echo "Guild Wars 2 Dev Tracker | https://en-forum.guildwars2.com/discover/6.xml"

go run . users
go run . feeds
