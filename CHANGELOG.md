# Changelog

## [2.0.0] - 2026.09.11

### Added
- Modals now include `type` field
- Modals now include option to drop files instead of asking for a reply with image

### Changed
- Rewritten in TypeScript using discordjs since discordgo was missing latest features from Discord API

### Fixed

## [1.2.0] - 2026.08.16

### Added

### Changed
- Database was changed from Supabase (Postgresql) -> Turso (libsql/sqlite)
- Removed docker-compose since database server no longer required for dev work
- Using database/sql instead of another dep for db connection

### Fixed

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
