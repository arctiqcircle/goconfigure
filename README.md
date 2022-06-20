# <img src="favicon.png" height="64"/> GoConfigure

> A simple SSH configuration deployment tool for the command-line.

## Operation

GoConfigure can be run at the command-line with `goconfigure`. Running `goconfigure` without
any arguments will launch an interactive session where target devices and commands can be
entered manually. Interactive mode does **not** support configuration templating.

Optional arguments include `-i inventory_filename` and `-t template_filename`. The template
should be defined in a plain-text document and the inventory filename should be defined in a
YAML or CSV formatted document. YAML should be supplied  according to the
following schema;

```yaml
---
default_username: User@domain.com
default_password: Password1
hosts:
  - hostname: server1.domain.com
    username: User@domain.com
    password: Password1
    data:
      Sample: hello server1!
  - hostname: server2.domain.com
    username: User@domain.com
    password: Password1
    data:
      Sample: hello server2!
```

*Note that `data`, `username`, and `password` are optional Host fields*.

The fields defined in `data` will be available to the
template during rendering. Fields required by the template must be defined in `data`. Field names
should be uppercase to be accessible from within the template.

It is not necessary to provide host specific credentials to the `hosts` definitions when
`default_username` or `default_password` are supplied.

CSV documents should follow the following schema. *Note that CSV does not support
`default_username` or `default_password`*

```csv
hostname,username,password,Sample
server1.domain.com,User@domain.com,Password1,"hello server1!"
server2.domain.com,User@domain.com,Password1,"hello server2!"
```

## Roadmap

The following is a roadmap of features to be added. Completing all the following will
get us to release v1.0.0.

- ~~Basic user/pass authentication~~.
- Key based authentication.
- ~~Implement YAML Inventories~~.
- ~~Implement CSV Inventories~~.
- Implement interactive operation.
- Index responses by command.
- Add support for recurring attempts.