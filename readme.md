# Mashiron-V ✌

[![Docker Image CI](https://github.com/tpc3/Mashiron-V/actions/workflows/docker-image.yml/badge.svg)](https://github.com/tpc3/Mashiron-V/actions/workflows/docker-image.yml)
[![Go](https://github.com/tpc3/Mashiron-V/actions/workflows/go.yml/badge.svg)](https://github.com/tpc3/Mashiron-V/actions/workflows/go.yml)

<div style="text-align: center;">
    <img src="./mashiron.png" width="100">
</div>

Advanced auto-reply bot, rewrite of the Mashiron-Fly with additional future.

## Requirement

* Any computer that runs Golang - hopefully.
  * Linux is our main environment for both dev and prod.
  * We don't check for other environments such as Windows, Mac, other architectures other then x86-64 like arm64, mips...
* Discord bot account (why not?)
* Time for deploy and maintain ~~bots~~ daughters
* Be ready for being alpha-tester in prod env

## Deployment

1. Install docker(if you want to)
    * ArchLinux: `pacman -Syu --needed docker`
1. [Download config](https://raw.githubusercontent.com/tpc3/Mashiron-V/master/config.yaml)
    * Adjust config for your env!
1. `docker run --rm -it -v $(PWD)/config.yaml:/Mashiron-V/config.yaml ghcr.io/tpc3/mashiron-v`
    * Simply download binary if you do not want to use docker
1. Profit

### Persistance

If you want to save data, bind these dirs:

* `/Mashiron-V/data`
* `/Mashiron-V/config`

## Migrate from "Fly"

Simply move `data` dir from Fly to V.  
V heavily rewrites data from Fly so be careful to always take some backups!

If you want to add "Fly"-style yaml from `add` command, use `--flex`.

## History

Few years ago, in 2019 if my memory is correct...  
At that time, we used LINE, the chat platform mainly used in Japan.  
On LINE, I made the very first auto-reply chat called "Mashiron" with Python.  
In that group, several other bot was working and she was one of those.

After some time, we said goodbye to LINE and switched to Discord.  
In discord, I also made the bot also called "Mashiron" (also written with Python).  
That was totally different from current Mashiron, since it was the control bot for Minecraft server and some systemd service daemon (really!).  
But you know, controlling minecraft server in some sort of web-ui is way easier for me and user, so that bot was deprecated super quickly(like sub-millisecond).  
For this reason(and this is even not the auto-reply bot), we usually don't count it as Mashiron(and I think almost all members don't remember that).

So what now? Well, Auto-reply bot. Again.  

On the discord server, we used dyno bot for auto-reply thing, but since we're small tech people group, that was not enough in two meanings:

* Not only admins, but we want all users able to "hack"[^1] discord server.
* We also want some programming future for week-end "hacking"[^1].

[^1]: Small coding for fun and improvement of the community.

According to these needs, I made a bot with same name, "Mashiron" with discord.py.  
It was the programmable auto-reply bot using bash.  
We really loved her, so rulesets are so huge that Python(in other words: my not optimized code) can't keep up with it.

I also loved loved her, so I made a bot called "Mashiron-Go!".  
That was the module-based bot written with golang, and able to connect to LINE and Discord.  

```

Discord => |                   | => Ping module
           | => Core binary => | => sh module
LINE    => |                   | => Weather module

```

It also has several modules like splatoon2 stage module.

...But as you can see, this mechanism is kinda really complicated and difficult to maintain.  
Also, writing readable bash code is difficult compared to another high-level languages.  
So I made a bot called "Mashiron-Fly", the auto-reply bot written by rust.  
It can simply define rules by writing tiny yaml(and javascript eval with quickjs if you want to do some coding).

```yaml
"Hi!": "How are you?"
advanced:
  trigger:
    - "^trigger$"
  return:
    - "Hello world!"
  js: "'hello, ' + scriptArgs[1] + '!'"
```

But this time, my coding was the problem.  
I came up this idea while playing Apex Legends so code was huge mess of spaghetti.  
There's only two rust coders in our community so maintenance was also huge problems.

So I re-written the bot with not-so-changed code quality in Golang, called "Mashiron-V".  
This is the fifth auto-reply bot, and our community is about to ready for five-years anniv!  
That's why I set the codename to "V", with her cute victory sign - ✌
