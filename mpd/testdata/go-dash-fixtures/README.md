# MPD fixtures from go-dash

The MPDs in this directory are copied from [go-dash](github.com/zencoder/go-dash).
They are unchanged except for

* `events.mpd` id had non-allowed type string. Changed to int.
* Durations have been simplified like: PT65.063S -> PT1M5.063S, PT2.000S -> PT2S, PT120S -> PT2M
