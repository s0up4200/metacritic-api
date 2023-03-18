# Metacritic Album API

This API scrapes album information from [Metacritic](https://www.metacritic.com/) for both [upcoming releases](https://www.metacritic.com/browse/albums/release-date/coming-soon/) and [new releases](https://www.metacritic.com/browse/albums/release-date/new-releases/date), and returns the data in JSON format. The API gathers artist names and album titles, updates the information every 24 hours, and caches the results to improve performance.

## Features

- Scrapes album information from Metacritic, specifically upcoming and new releases
- Returns artist names and album titles in JSON format
- Caches data and updates every 24 hours to minimize resource usage
- Prevents race conditions using mutex locks
