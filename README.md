# docker-arccheck

Docker container including tools for Adaptec RAID controllers

tools included :

- arcconf
- check_adaptec_raid from [thomas-krenn](https://github.com/thomas-krenn/check_adaptec_raid)
- arccheck : go monitoring daemon that run check_adaptec_raid and send telegram notification on state transition

## Usage :

### arcconf
`docker run --rm --privileged -it akit042/docker-arccheck arcconf GETCONFIG`

### check_adaptec_raid
`docker run --rm --privileged -it akit042/docker-arccheck check_adaptec_raid -Tw 40 -Tc 50 -LD 0,1 -PD 1 -z 0`

### arccheck
`docker run --rm --privileged -it akit042/docker-arccheck arccheck -telegramtoken XXXXXXX:XXXXXXXXXXXXXXX -telegramid YYYYYYYY -commandargs '-Tw 40 -Tc 50 -LD 0,1 -PD 1 -z 0' -poolinginterval 10`

```
Usage of arccheck:
  -commandargs string
        Arguments to use for check_adaptec_raid (or use env variable : COMMAND_ARGS)
  -d    debug mode
  -poolinginterval int
        Pooling Interval (or use env variable : POOLING_INTERVAL) (default 30 minutes)
  -telegramid int
        To find an id, please contact @myidbot on telegram (or use env variable : TELEGRAM_ID)
  -telegramtoken string
        To create a bot, please contact @BotFather on telegram (or use env variable : TELEGRAM_TOKEN)
  -v    Print build id
```
