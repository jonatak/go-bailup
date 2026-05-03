# Changelog

## [1.2.3](https://github.com/jonatak/baillconnect-to-mqtt/compare/v1.2.2...v1.2.3) (2026-05-03)


### Bug Fixes

* add broader permission on apparmor ([848b78a](https://github.com/jonatak/baillconnect-to-mqtt/commit/848b78ac7c5a7acd08499d6a0e1efa8d2f44545c))

## [1.2.2](https://github.com/jonatak/baillconnect-to-mqtt/compare/v1.2.1...v1.2.2) (2026-05-03)


### Bug Fixes

* another apparmor fix ([998161c](https://github.com/jonatak/baillconnect-to-mqtt/commit/998161c1af10693400d8a4295ff6a439add8bf28))

## [1.2.1](https://github.com/jonatak/baillconnect-to-mqtt/compare/v1.2.0...v1.2.1) (2026-05-03)


### Bug Fixes

* update apparmor to allow base homeassistant permission ([518f3e1](https://github.com/jonatak/baillconnect-to-mqtt/commit/518f3e108d4ae5bcab63787f6dce3b6df6219154))

## [1.2.0](https://github.com/jonatak/baillconnect-to-mqtt/compare/v1.1.0...v1.2.0) (2026-05-03)


### Features

* add apparmor ([dbeede8](https://github.com/jonatak/baillconnect-to-mqtt/commit/dbeede86c9c6b5095a5d39abca710593565d5825))
* make sure mqtt handler can send new cmd in a non blocking maner ([fbc1333](https://github.com/jonatak/baillconnect-to-mqtt/commit/fbc13332dbef81024032db2b6a2c820846e03ca7))
* mutualise installation script with docker build ([5153b14](https://github.com/jonatak/baillconnect-to-mqtt/commit/5153b14a8bcb9a1e756ecb3be5d667e5cde47164))


### Bug Fixes

* exit with error code in case of error ([a0bcb93](https://github.com/jonatak/baillconnect-to-mqtt/commit/a0bcb9313ea82735ad1b0a3de77ffa30c860416b))

## [1.1.0](https://github.com/jonatak/baillconnect-to-mqtt/compare/v1.0.1...v1.1.0) (2026-05-03)


### Features

* add battery monitoring in mqtt ([b845d7e](https://github.com/jonatak/baillconnect-to-mqtt/commit/b845d7eff43d34eaea0b832f58ea2d2094f03ca3))
* add battery state in domain ([c6bca8b](https://github.com/jonatak/baillconnect-to-mqtt/commit/c6bca8b7eaa586cefe5593d9677537bcd47e087e))
* improve cicd ([8a1bce5](https://github.com/jonatak/baillconnect-to-mqtt/commit/8a1bce5fd7eb70b81dbb54951841bc937af91494))


### Bug Fixes

* fix flaky viper env test ([f080875](https://github.com/jonatak/baillconnect-to-mqtt/commit/f0808753d2f25765328f4779ae3a036d23e91690))

## [1.0.1](https://github.com/jonatak/baillconnect-to-mqtt/compare/v1.0.0...v1.0.1) (2026-05-01)


### Bug Fixes

* correct changelog order in cicd ([8cda8a2](https://github.com/jonatak/baillconnect-to-mqtt/commit/8cda8a2097c08d3ecf4a8f802e45ddcb0047c0ae))

## [1.0.0](https://github.com/jonatak/baillconnect-to-mqtt/compare/v0.2.0...v1.0.0) (2026-05-01)


### ⚠ BREAKING CHANGES

* first release of HAOS plugin ([#7](https://github.com/jonatak/baillconnect-to-mqtt/issues/7))

### Features

* add small install script ([ee5cb5e](https://github.com/jonatak/baillconnect-to-mqtt/commit/ee5cb5ec92f0167deb9e6efc2e050961e670c926))
* first release of HAOS plugin ([#7](https://github.com/jonatak/baillconnect-to-mqtt/issues/7)) ([e370239](https://github.com/jonatak/baillconnect-to-mqtt/commit/e370239c1566f293c311b12b5aaa2da1f82292ae))


### Bug Fixes

* typo in cicd ([4f6db79](https://github.com/jonatak/baillconnect-to-mqtt/commit/4f6db79b4af95c003c29e6c295b8a745169a83b2))

## [0.2.0](https://github.com/jonatak/baillconnect-to-mqtt/compare/v0.1.3...v0.2.0) (2026-04-29)


### Features

* implement ci with releaser ([197bfd2](https://github.com/jonatak/baillconnect-to-mqtt/commit/197bfd2e15ac31c49124b170e62ec69796a729e2))


### Bug Fixes

* correct linter issues ([ac1707f](https://github.com/jonatak/baillconnect-to-mqtt/commit/ac1707fcac377006cf8e35574c5ebbcf3e029122))
