# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased
### Added
- /about endpoint with build information
- Version of the webapp in the webapp's header
- ForwarderVersion and IndexerVersion in the indexed fields
### Changed
- The `c.analyzer` and `go.analyzer` options now take the gdb or delve commands to execute on the coredump, instead of full shell commands (#27)
- Design & Interface overhaul:
	- Better table readability and usage by changing padding and making the whole row clickable
	- Rework the searchbar to remove useless options and optimize the sorting
	- Reorganization of the detail view and table heading for better readability
- The sort and sort order parameters are now separate options because it's simpler
- Limit the sort options to date and hostname because the executable one isn't working as intended anymore
### Fixed
- Case of the fields for sorting is back to lowercase (#26)

## [v0.6.1](https://github.com/elwinar/rcoredump/releases/tag/v0.6.1) - 2020-03-23
### Internal
- Cleanup `go.mod` dependencies for unused packages

## [v0.6.0](https://github.com/elwinar/rcoredump/releases/tag/v0.6.0) - 2020-03-22
### Added
- `metadata` options accept a list of key-value pairs to send alongside the coredump to the indexer
- Travis-CI integration is setup [here](https://travis-ci.org/github/elwinar/rcoredump)
### Changed
- the `dir` option becomes `data-dir`
### Removed
- `bleve.path` configuration option: it is now deduced from the `data-dir` option