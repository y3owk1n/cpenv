# Changelog

## [1.13.0](https://github.com/y3owk1n/cpenv/compare/v1.12.1...v1.13.0) (2024-12-08)


### Features

* add short aliases for commands ([#100](https://github.com/y3owk1n/cpenv/issues/100)) ([9595a16](https://github.com/y3owk1n/cpenv/commit/9595a165092f5599c0690d9a9381a4ad6b574587))
* better description and title for prompts ([#105](https://github.com/y3owk1n/cpenv/issues/105)) ([b4148f9](https://github.com/y3owk1n/cpenv/commit/b4148f9d2e6538f6f14a9f22c5dadf1e657794ff))
* change alias from vc to v ([#102](https://github.com/y3owk1n/cpenv/issues/102)) ([b1d5534](https://github.com/y3owk1n/cpenv/commit/b1d553447f0eb4a072570ec93e1fea19661b6a81))
* remove init function ([#98](https://github.com/y3owk1n/cpenv/issues/98)) ([b321578](https://github.com/y3owk1n/cpenv/commit/b32157843d23ea5d2c987816a129c73d3dd8d60a))
* use base theme instead of colored theme ([#103](https://github.com/y3owk1n/cpenv/issues/103)) ([17d01e3](https://github.com/y3owk1n/cpenv/commit/17d01e3823cfde2ed41b43e9fd9ed8cf4342ff0c))
* use built-in ErrUserAborted instead of == check ([#104](https://github.com/y3owk1n/cpenv/issues/104)) ([244ccad](https://github.com/y3owk1n/cpenv/commit/244ccad70b6ded73c15b5ae0cdfb35bcd4a28884))


### Bug Fixes

* exit if not confirm cwd for backup ([#101](https://github.com/y3owk1n/cpenv/issues/101)) ([5a80c67](https://github.com/y3owk1n/cpenv/commit/5a80c67dd0439024b5d274394289586cf1976c3a))

## [1.12.1](https://github.com/y3owk1n/cpenv/compare/v1.12.0...v1.12.1) (2024-12-07)


### Bug Fixes

* log "user aborted" as debug level ([#97](https://github.com/y3owk1n/cpenv/issues/97)) ([98db947](https://github.com/y3owk1n/cpenv/commit/98db947e0ee688549d350acb3af688fd74d6de9a))
* remove debugLevel for prod ([#95](https://github.com/y3owk1n/cpenv/issues/95)) ([b2a3c1d](https://github.com/y3owk1n/cpenv/commit/b2a3c1dab2d732f1a0dd27409f013645c59de3c5))

## [1.12.0](https://github.com/y3owk1n/cpenv/compare/v1.11.2...v1.12.0) (2024-12-07)


### Features

* refactor cmd commands for better handling ([#90](https://github.com/y3owk1n/cpenv/issues/90)) ([734eed0](https://github.com/y3owk1n/cpenv/commit/734eed0da2bec82d41ecb82119c34a28e7316eaf))
* refactor core and utils for better maintainence ([#92](https://github.com/y3owk1n/cpenv/issues/92)) ([97a26d8](https://github.com/y3owk1n/cpenv/commit/97a26d8b909a816fb832b320540487ad85d00112))


### Bug Fixes

* index not found, should append instead ([#93](https://github.com/y3owk1n/cpenv/issues/93)) ([3f108ee](https://github.com/y3owk1n/cpenv/commit/3f108ee19436e2652ee2282f346f9e5b4e9f4042))
* make sure to exit the os if no env configuration ([#94](https://github.com/y3owk1n/cpenv/issues/94)) ([c7c24d4](https://github.com/y3owk1n/cpenv/commit/c7c24d472f853a2489a4583724be8e838dcea355))

## [1.11.2](https://github.com/y3owk1n/cpenv/compare/v1.11.1...v1.11.2) (2024-12-07)


### Bug Fixes

* use the right namespace to build for version to show ([#88](https://github.com/y3owk1n/cpenv/issues/88)) ([49e29ce](https://github.com/y3owk1n/cpenv/commit/49e29ce2f8042b0f7f266ba35998e8c61188e5d5))

## [1.11.1](https://github.com/y3owk1n/cpenv/compare/v1.11.0...v1.11.1) (2024-12-07)


### Bug Fixes

* refactor messages icon ([#85](https://github.com/y3owk1n/cpenv/issues/85)) ([51ce119](https://github.com/y3owk1n/cpenv/commit/51ce1198f95dd75faf978939a20567cc8fd4c38a))
* remove unnecssary println for copy & backup completed ([#84](https://github.com/y3owk1n/cpenv/issues/84)) ([295b08f](https://github.com/y3owk1n/cpenv/commit/295b08f65f11a416cc9b954e03bb7b102484c6fc))
* update namespace to url style ([#87](https://github.com/y3owk1n/cpenv/issues/87)) ([81edfaa](https://github.com/y3owk1n/cpenv/commit/81edfaad8a712228e4a60fbdde67ff1d692027ef))
* use form and inline catppuccin theme ([#82](https://github.com/y3owk1n/cpenv/issues/82)) ([8423869](https://github.com/y3owk1n/cpenv/commit/8423869cfd289f4eaa3fef5d99debc0acdf64643))

## [1.11.0](https://github.com/y3owk1n/cpenv/compare/v1.10.0...v1.11.0) (2024-12-07)


### Features

* add catppuccin theme for huh forms ([#81](https://github.com/y3owk1n/cpenv/issues/81)) ([5624a54](https://github.com/y3owk1n/cpenv/commit/5624a54ac21c1da45fdc9f03cd0a074e91840a53))
* bump go deps ([#79](https://github.com/y3owk1n/cpenv/issues/79)) ([c41215b](https://github.com/y3owk1n/cpenv/commit/c41215b489cb39cea721552d4f4616e644c22eec))

## [1.10.0](https://github.com/y3owk1n/cpenv/compare/v1.9.0...v1.10.0) (2024-12-07)


### Features

* move version command to flag instead ([#77](https://github.com/y3owk1n/cpenv/issues/77)) ([5ba7d78](https://github.com/y3owk1n/cpenv/commit/5ba7d78e73239c3728abde6387553ecfc05ef1be))

## [1.9.0](https://github.com/y3owk1n/cpenv/compare/v1.8.3...v1.9.0) (2024-12-07)


### Features

* add vault command to open in finder ([#75](https://github.com/y3owk1n/cpenv/issues/75)) ([adf17b2](https://github.com/y3owk1n/cpenv/commit/adf17b23e296fc6026687085b2c202cdaf4af226))


### Bug Fixes

* gracefully exit huh form when user aborted ([#73](https://github.com/y3owk1n/cpenv/issues/73)) ([e5eeafd](https://github.com/y3owk1n/cpenv/commit/e5eeafdf335ad0047fc9af7d5b9d38442ea419e7))

## [1.8.3](https://github.com/y3owk1n/cpenv/compare/v1.8.2...v1.8.3) (2024-12-03)


### Bug Fixes

* wrong build for arm64!!!! ([#71](https://github.com/y3owk1n/cpenv/issues/71)) ([4372748](https://github.com/y3owk1n/cpenv/commit/4372748f2ebae4b0194ef4278b8355a279f040af))

## [1.8.2](https://github.com/y3owk1n/cpenv/compare/v1.8.1...v1.8.2) (2024-12-03)


### Bug Fixes

* set log level to default, debug is for local dev ([#69](https://github.com/y3owk1n/cpenv/issues/69)) ([43f9314](https://github.com/y3owk1n/cpenv/commit/43f931426257700255c2ffbd09a5b30a8a8d45a7))

## [1.8.1](https://github.com/y3owk1n/cpenv/compare/v1.8.0...v1.8.1) (2024-12-03)


### Bug Fixes

* wrong build for arm64 ([#67](https://github.com/y3owk1n/cpenv/issues/67)) ([933f132](https://github.com/y3owk1n/cpenv/commit/933f132adf43b74f6664ebd69ee1a74169e7469e))

## [1.8.0](https://github.com/y3owk1n/cpenv/compare/v1.7.0...v1.8.0) (2024-12-03)


### Features

* migrate to golang ([#63](https://github.com/y3owk1n/cpenv/issues/63)) ([7eed464](https://github.com/y3owk1n/cpenv/commit/7eed464ac4366443e3e74f557d63bb36f438f8e0))

## [1.7.0](https://github.com/y3owk1n/cpenv/compare/v1.6.1...v1.7.0) (2024-12-03)


### Features

* add setup command ([#61](https://github.com/y3owk1n/cpenv/issues/61)) ([f5f2035](https://github.com/y3owk1n/cpenv/commit/f5f203505fedc807761b2ce287704431e5688cd7))


### Bug Fixes

* revert root call for interactive actions ([#58](https://github.com/y3owk1n/cpenv/issues/58)) ([9f14d72](https://github.com/y3owk1n/cpenv/commit/9f14d72584f717664a1fd719bcd8023e63ff11b7))

## [1.6.1](https://github.com/y3owk1n/cpenv/compare/v1.6.0...v1.6.1) (2024-12-02)


### Bug Fixes

* do not log silly thing on sigTerm ([#56](https://github.com/y3owk1n/cpenv/issues/56)) ([a965974](https://github.com/y3owk1n/cpenv/commit/a9659745832348a1f0cd396450b3b8c3685a641d))

## [1.6.0](https://github.com/y3owk1n/cpenv/compare/v1.5.0...v1.6.0) (2024-12-02)


### Features

* add subcommand for copy & backup ([#52](https://github.com/y3owk1n/cpenv/issues/52)) ([d724bfd](https://github.com/y3owk1n/cpenv/commit/d724bfd0bc2ab7bbe1949eb46473104c9d10b7b4))


### Bug Fixes

* make sure to handle exit event properly ([#54](https://github.com/y3owk1n/cpenv/issues/54)) ([50ede4f](https://github.com/y3owk1n/cpenv/commit/50ede4f378e5637efa6508aa818011049a0b050b))

## [1.5.0](https://github.com/y3owk1n/cpenv/compare/v1.4.2...v1.5.0) (2024-12-02)


### Features

* handle exit error globally ([#49](https://github.com/y3owk1n/cpenv/issues/49)) ([8018475](https://github.com/y3owk1n/cpenv/commit/8018475aea283490961aee17eac2051e591df18a))
* use exact bun version to build ([#51](https://github.com/y3owk1n/cpenv/issues/51)) ([72ed67b](https://github.com/y3owk1n/cpenv/commit/72ed67bbcf75647b420665ef05310d2467f462e3))

## [1.4.2](https://github.com/y3owk1n/cpenv/compare/v1.4.1...v1.4.2) (2024-12-01)


### Bug Fixes

* version bug causes binary to not run ([#45](https://github.com/y3owk1n/cpenv/issues/45)) ([3d6bc98](https://github.com/y3owk1n/cpenv/commit/3d6bc98c40c2321620aac6085e27be94a6df5b14))

## [1.4.1](https://github.com/y3owk1n/cpenv/compare/v1.4.0...v1.4.1) (2024-12-01)


### Bug Fixes

* make sure to use cjs instead of js ([#41](https://github.com/y3owk1n/cpenv/issues/41)) ([500e716](https://github.com/y3owk1n/cpenv/commit/500e7160e4d0636f504874a9e8ebb7a6a23918e0))

## [1.4.0](https://github.com/y3owk1n/cpenv/compare/v1.3.1...v1.4.0) (2024-12-01)


### Features

* add action to test the build ([#39](https://github.com/y3owk1n/cpenv/issues/39)) ([82472f5](https://github.com/y3owk1n/cpenv/commit/82472f5730036cc2b173c8219186e103359ab3a4))

## [1.3.1](https://github.com/y3owk1n/cpenv/compare/v1.3.0...v1.3.1) (2024-12-01)


### Bug Fixes

* add release id to gh action ([#37](https://github.com/y3owk1n/cpenv/issues/37)) ([13a09f4](https://github.com/y3owk1n/cpenv/commit/13a09f4dcc6d9f21faee17bf235ae2329a48f861))

## [1.3.0](https://github.com/y3owk1n/cpenv/compare/v1.2.0...v1.3.0) (2024-12-01)


### Features

* add release please config ([#32](https://github.com/y3owk1n/cpenv/issues/32)) ([deb3cdb](https://github.com/y3owk1n/cpenv/commit/deb3cdb12d27b9d8207539487a6cfee5ec39b0b3))
* update gh actions for auto release with binary uploads ([#34](https://github.com/y3owk1n/cpenv/issues/34)) ([c0c9c79](https://github.com/y3owk1n/cpenv/commit/c0c9c79a12b1fa775715bacd516538f50618341e))


### Bug Fixes

* make sure bun setup dont run during pr ([#36](https://github.com/y3owk1n/cpenv/issues/36)) ([3b8049b](https://github.com/y3owk1n/cpenv/commit/3b8049b0604aebf41c445f314a2c6a00dd882a48))
