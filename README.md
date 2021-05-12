# deckd

## Disclaimer

**¡ This project is highly experimental !**

Although it might be a reliable piece of software at some point, at the moment
it is developed with educational intentions in mind. In other words: it's my personal
playground for learning golang.

## Intro

Deckd is a playback engine meant to be used as the core for a DJ Application.
Its purpose is to support the least number of operations that typical features of
a DJ App can be broken down to. This keeps complexity low but(hopefully) makes
playback as reliable as possible.

At the moment `jack` is the only output backend supported. It is developed using
`Pipewires` implementation of `jack`.
