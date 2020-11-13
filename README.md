
# Mattermost as code

As our Mattermost server regularly resets all configured sidebar categories,
I wrote this tool to be able to export/import them automatically again.


## Installation

You can either compile it yourself, e.g. via `go get` or download it from
the [releases](https://github.com/sbueringer/mattermost-as-code/releases) page.

## Usage

Export the current configuration for a specific team:

```bash
export MATTERMOST_USERNAME=
export MATTERMOST_PASSWORD=
export MATTERMOST_URL=
```

```bash
mac export --team <team-name> > ./export.yaml
```

Import a configuration for a specific team:

```bash
mac import --input ./export.yaml
```


## Configuration file format

```yaml
teams:
- name: <team-name>
  categories:
  # Favorites
  - channels:
    - <channel-name>
    - <channel-name>
    name: Favorites
    sorting: manual
    type: favorites

  # Custom Category 1
  - channels:
    - <channel-name>
    - <channel-name>
    name: <category-name>
    sorting: manual
    type: custom
  
  # Custom Category 2
  - channels:
    - <channel-name>
    - <channel-name>
    name: <category-name>
    sorting: manual
    type: custom
```

## Build

It can also be build via [mage](https://github.com/magefile/mage):

```bash
mage -v buildALl
ls -la ./dist

mage -v lint
mage -v format
mage -v coverage
mage -v clean
```

To run the Github pipeline locally just use [act](https://github.com/nektos/act):

```bash
act -l
act
```
