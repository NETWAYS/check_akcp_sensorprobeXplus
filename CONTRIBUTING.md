# Contributing
Before starting your work on this module, you should [fork the project] to your GitHub account. This allows you to
freely experiment with your changes. When your changes are complete, submit a [pull request]. All pull requests will be
reviewed and merged if they suit some general guidelines:

* Changes are located in a topic branch
* For new functionality, proper tests are written
* Changes should not solve certain problems on special environments

## Branches
Choosing a proper name for a branch helps us identify its purpose and possibly find an associated bug or feature.
Generally a branch name should include a topic such as `fix` or `feature` followed by a description and an issue number
if applicable. Branches should have only changes relevant to a specific issue.

```
git checkout -b fix/service-template-typo-1234
git checkout -b feature/config-handling-1235
git checkout -b doc/fix-typo-1236
```

## Building

Download and install Go according to your operating system.

```
git clone https://github.com/NETWAYS/check_akcp_sensorprobeXplus
cd check_akcp_sensorprobeXplus
go build
```
