# Changelog

## [1.1.1] - 2026.07.24

### Added
- Dockerfile for containerization

### Changed

### Fixed

## [1.1.0] - 2026.07.20

### Added

### Changed
- Users can now change `name` field in the edit modal (bot will check if name already exists in database before submitting)
- Users can now choose to either `keep`, `change` or `delete` image associated to a card.

### Fixed
- Fixed a visual bug where some replies from the bot looked like `Card for {name} ** saved` instead of `Card for ** {name} ** saved`

## [1.0.0] - 2026.07.20

### Added
- Initial release
- Basic CRUD operations for lore cards
- For `create` and `update`, the bot uses forms for text data and a goroutine is spun up awaiting an attachement

### Changed

### Fixed
