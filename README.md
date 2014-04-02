# quasar

quasar is IRC bot to help with various (and mostly tedious) tasks that we
encountered while extensively using IRC for communication at
[Kir-Dev](http://kir-dev.sch.bme.hu).

## Main features

* first and foremost: logging
* a nice web interface for catching up with logs
* github integration for linking issues
* github integration for notifications (new issue, closed issue, new comment, etc)
* inlining tweets
* youtube integration: title and description of videos
* reporting

A comprehensive list is available on the [Kir-Dev trello board](https://trello.com/b/HN5AgGe8/kir-dev-board) (in Hungarian).

## Overview

quasar is actually two separate applications. The web interface is located under
the `quasar-web` directory. It is a [Node.js](http://nodejs.org) application. The bot
itself is (surprisingly) located under the `quasar-bot` directory. It is a
[go](http://golang.org) application.

[PostgreSQL](http://www.postgresql.org/) used for storage.

## Installing

### Prerequisites

1. install the latest [PostgreSQL](http://www.postgresql.org/download/)
2. install [go1.2+](http://golang.org/doc/install#download)
3. install Node.js

### Setting up the database

Assuming a running and functional postgresql server.

    $ sudo su - postgres
    $ createuser -l -E -P -R -S -D quasar
    $ createdb -O quasar -E utf8 quasar
    $ psql -U quasar -d quasar -h localhost -f /path/to/quasar/scripts/sql/schema.sql

### Setting up your go environment

Assuming you have go installed, create your workspace.

    $ mkdir -p /path/to/your/workspace/src

The src folder at the end is mandatory.

### Get the code

    $ cd /path/to/your/workspace/src
    $ git clone https://github.com/kir-dev/quasar.git quasar

### Setup your config file

    $ cp config/config.json.dist config/config.json
    $ vim config/config.json

For detailed information about the configuration look at the
[configuration](#configuration) section.

### Build & run quasar-bot

We are using [godep](https://github.com/tools/godep) for managing dependencies,
so you must have `godep` installed in your `PATH`.

To build the bot itself just run

    $ make bot

It creates a new executable named `quasar`. To run it simply:

    $ ./quasar

### Setting up the node.js environment

TODO

## Configuration

TODO

## Committing

When commiting go code **always** use the `go fmt` tool first. Possibly one could
set up a pre-commit git-hook to automate this.

Or you can do it manually:

    $ go fmt ./...
    # or
    $ make fmt

## Adding a new dependency to quasar-bot

We are using [godep](https://github.com/tools/godep) for managing dependencies,
so you must have `godep` installed in your `PATH`.

Use the [godep workflow](https://github.com/tools/godep#add-or-update-a-dependency)
and commit the vendorized package source.