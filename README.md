# <img src="favicon.png" height="64"/> GoConfigure

> A simple SSH configuration deployment tool for the command-line.

# Operation

GoConfigure can be run at the command-line with `goconfigure`. Running `goconfigure` without
any arguments will launch an interactive session where target devices and commands can be
entered manually. Interactive mode does **not** support configuration templating.

Optional arguments include `-i inventory_filename` and `-t template_filename`. The template
should be defined in a plain-text document and the inventory filename should be defined in a
YAML formatted document according to the following schema;
```yaml
---
server-1:
  hostname: s1.yourdomain.com
  username: username
  password: password
  data:
    ip_address: 192.168.0.1
server-2:
  hostname: s2.yourdomain.com
  username: username
  password: password
```

*Note that `data` is an optional field*. The fields defined in `data` will be available to the
template during rendering. Fields required by the template must be defined in `data`. Field names
should be uppercase to be accessible from within the template.