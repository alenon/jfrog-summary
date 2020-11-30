# Artifactory Summary JFrog CLI plugin
Artifactory summary live visualisation (currently supported storage summary only).

## Installation with JFrog CLI
Installing the latest version:

`$ jfrog plugin install rt-summary`

Installing a specific version:

`$ jfrog plugin install rt-summary@version`

Uninstalling a plugin

`$ jfrog plugin uninstall rt-summary`

## Usage
### Commands
* storage - Artifactory storage summary

    - Usage: `jfrog rt-summary storage [command options]`

    - Options:
        - **server-id** - Artifactory server ID configured using the config command *[Optional]*
        - **refresh-rate** - summary refresh rate in seconds *[Default: 2]* 
        - **recalculate-rate** - storage summary recalculation rate in seconds. If 0 recalculation will not be triggered *[Default: 0]*
        - **max-results** - maximal amount of shown results *[Default: 10]*
    - Example:
    ```
  $ jfrog rt-summary st
  
    Last updated at: Sun, 22 Nov 2020 10:11:27 IST
    Last recalculated at: Sun, 22 Nov 2020 10:11:27 IST
    
    ┼────────────────────────┼───────┼──────────────┼───────────────┼──────────────┼──────────────┼
    │       REPOSITORY       │ TYPE  │ PACKAGE TYPE │  FILES COUNT  │  USED SPACE  │  PERCENTAGE  │
    ┼────────────────────────┼───────┼──────────────┼───────────────┼──────────────┼──────────────┼
    │ example-repo-local     │ LOCAL │ Generic      │ 1             │ 10.64 MB     │ 97.45%       │
    │ jfrog-support-bundle   │ NA    │ NA           │ 7             │ 285.31 KB    │ 2.55%        │
    │ artifactory-build-info │ LOCAL │ BuildInfo    │ 0             │ 0 bytes      │ 0%           │
    │ auto-trashcan          │ NA    │ NA           │ 0             │ 0 bytes      │ 0%           │
    ┼────────────────────────┼───────┼──────────────┼───────────────┼──────────────┼──────────────┼
    │         TOTAL          │   -   │      -       │       8       │   10.92 MB   │      -       │
    ┼────────────────────────┼───────┼──────────────┼───────────────┼──────────────┼──────────────┼
  ```

## Release Notes
The release notes are available [here](RELEASE.md).
