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
