# dungeonfs

DungeonFS is a FUSE filesystem and dungeon crawling adventure game engine.  This is a work-in-progress, however, there is a general list of [planned features](#roadmap) worth checking out, as well as, a little bit of the motivation behind [why I started this project](#motivation).

[![demo](https://asciinema.org/a/110084.png)](https://asciinema.org/a/110084?autoplay=1)

# Install

```Shell
go get github.com/ChrisRx/dungeonfs/...
```

Note: The demo level is included in the `<project root>/examples/simplelevel` folder, which should be downloaded with the above command. In the future, I will have this included with the static binary so the demo level is accessible without need go installed on the system.

This has also shown not to work on Go version 1.3 with failures building the `golang.org/x/sys/unix` package.  I will update this further when I figure out the actual minimum version supported by dungeonfs.

# Usage

The command-line `mount` command is used to mount the FUSE filesystem.  It can be used to mount as a background process with the following:

```Shell
bin/dungeonfs mount <mountpoint> -d -a <asset folder>
```

Since this is running in the background it must be unmounted using the `unmount` command:

```Shell
bin/dungeonfs unmount <mountpoint>
```

# Roadmap

As mentioned above, this a work-in-progress and many features are either either missing or incomplete.  These are all features that I feel will probably be very important to having a good core game engine (at least that I have thought of so far).  If you have other ideas on what would make the engine better to use from a designer perspective, please don't hesitate to start an issue.

- NPC/combat system
  - file size will represent NPC/enemy health
  - talking with NPCs starts a talk repl to make that easier
- Implicit 2d directions for navigation
  - implicit but can be set
  - opens up the ability to have a viewable map that marks the player position via `cat .inventory/map`
  - will disambiguate between folder name and traversal (e.g. folder is 'north' but room name is unique to the room)
- Level editor using [termbox](https://github.com/nsf/termbox-go)
- A helper script that can be loaded with `source` at the beginning of the dungeon to make helpful aliases
  - ensure PS1 is set to only show the current folder
  - specific settings to sorting, hidden, etc for commands like `ls`
  - set environment variables that will be helpful like an `$EXIT` (although exit could also be a function possibly)
- Scripting language features
  - some Go stdlib exposed
    - fmt, strings, sync, time, etc
  - async functions (i.e. `go func() { ... }`)
  - some game-specific locks to ensure consistency of the game environment between concurrent scripts (currently called properties)
  - builtins that help with calculating a property only N times
    - e.g. i only want the room to have the `hidden` attribute property run until it is changed to `false`, then no longer run the property
  - currently `node` is the selector for the current node within a property, but this could change if need be
  - add selectors the for the names of over nodes for easy access to other node attributes
 
# Ideas

These are just some things I wanted to collect into a list even though I'm not sure if they will all be possible (or how they will work quite yet).  Really just didn't want to lose track of them so that they can be explored.

- Sound effects possible via aplay
- Using `chown` and `mv` in creative ways to buy/sell or steal
- Multi-player (issue #2)

# Motivation

I've always thought FUSE filesystems were neat and when I came upon the [https://github.com/bazil/fuse](https://github.com/bazil/fuse) package I wanted to make a trivial example to show what kind of crazy stuff could be abstracted into a filesystem.  Initially, the idea was just to make a simple dungeon with directories as rooms, but that ended up spiraling into also wanting to learn more about game programming (having never done much game programming).  Imade this as a learning project for myself, but as it became less about making a single example and more aspiring to be a game engine, I wanted to share it with anyone that might find it interesting to use and/or to develop the engine.  This project has been so fun and rewarding to work on, and my hope now is to make the game engine more complete so that others can try to use it to make their own novelty filesystem games.
