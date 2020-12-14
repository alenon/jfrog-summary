# Artifactory Summary JFrog CLI plugin
Artifactory summary live visualisation (currently supported storage summary only).

## Installation with JFrog CLI
Since this plugin is currently not included in [JFrog CLI Plugins Registry](https://github.com/jfrog/jfrog-cli-plugins-reg), it needs to be built and installed manually. Follow these steps to install and use this plugin with JFrog CLI.
1. Make sure JFrog CLI is installed on you machine by running ```jfrog```. If it is not installed, [install](https://jfrog.com/getcli/) it.
2. Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not exist already.
3. Clone this repository.
4. CD into the root directory of the cloned project.
5. Run ```go build``` to create the binary in the current directory.
6. Copy the binary into the ```~/.jfrog/plugins``` directory.

## Usage
### Commands
* storage - Artifactory storage summary

    - Usage: `jfrog rt-summary storage [command options]`

    - Options:
        - **server-id** - Artifactory server ID configured using the config command *[Optional]*
        - **live** - live summary update *[Default: false]*
        - **repo-list** - comma separated repositories list to show *[Default: all]*
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
