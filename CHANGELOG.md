# Changelog

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
