# boomtown (wildcatting 2)

## API

POST    /game/?weeks=<n>          - create a game
GET     /game/                    - list games
GET     /game/<id>/               - game summary
POST    /game/<id>/player/<name>/ - join/begin/survey/drill/sell
GET     /game/<id>/player/<name>/ - player state
