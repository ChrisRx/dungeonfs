# dungeonfs

A FUSE filesystem and dungeon crawling adventure game engine.

[![demo](https://asciinema.org/a/110084.png)](https://asciinema.org/a/110084?autoplay=1)


# Install

```Shell
make deps
make
```

# Usage

The command-line `mount` command is used to mount the FUSE filesystem.  It can be used to mount as a background process with the following:

```Shell
bin/dungeonfs mount <mountpoint> -d
```

Since this is running in the background it must be unmounted using the `unmount` command:

```Shell
bin/dungeonfs unmount <mountpoint>
```

# Ideas

- Sound effects possible via aplay
- NPCs files
- Using `chown` and `mv` in creative ways to buy/sell or steal
- NPC/enemy health and other attributes can be determined via file attributes like filesize
- Disambiguate between folder name and traverse (e.g. folder is 'north' but room name is unique to the room)
- Add `.inventory/map` that can allow the player to view the map.
  - Directions not made explicit can be enabled by converting the name to a number and then mod that number by the number of directions desired for the map (e.g. 4 for basic cardinal directions north,south,east,west)
- Level editor using something like [termbox](https://github.com/nsf/termbox-go).  Could also be used to make animations that run when doing certain actions (e.g. looting a chest).



# TODO

- Expose go stdlib to scripting language: fmt, strings, sync, time, etc
- Add some builtins for scripting language:
  - A simple global lock that must be acquired by the fs methods
- Shell script that can be sourced that provides the recommended settings for commands like `ls` and global variables like `$EXIT` (or possible a function) that helps to exit the file system.  This would have to be run with `source` to work.
