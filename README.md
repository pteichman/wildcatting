# Wildcatting

## East Texas, Gusher Age

*The easy-going rural life of East Texas changed drastically with the
discovery of oil in 1930 and 1931 – years of hardship, scorn, luck and wealth which brought people, ideas, institutions and national attention to East Texas.*

![Gusher](wildcatting.jpg)

## API

    POST    /game/                     - create -> gameID
    GET     /game/<id>/                - game status
    POST    /game/<id>/                - join -> playerID
    POST    /game/<id>/player/<id>/    - start/survey/drill/sell -> player view
    GET     /game/<id>/player/<id>/    - player view

## Bootstrap

Create a game and join a player, as there is no UI for this stuff yet:

    $ curl -X POST http://localhost:8888/game/ && curl -X POST http://localhost:8888/game/0/ -d '"bob"'
