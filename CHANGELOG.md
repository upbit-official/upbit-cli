# Changelog

## [1.0.1]
- Require Go SDK v1.0.1 so `--environment id` uses `https://api.upbit.co.id/`; retain explicit URL overrides.
- Redact sensitive request and response headers regardless of case, duplicate keys, or multiple values.
- Fix decimal percentage calculations in the take-profit/stop-loss shell examples while preserving the final price floor.
- Refresh generated examples, including `new_volume` for cancel-and-create orders and omitted optional parameters.
- Correct local Go SDK development instructions.

## [1.0.0]
- Add Pocket command (`pkg/cmd/pocket.go`): 7 new REST endpoints
- Add missing `unit` field to minute candle response

## [0.9.3]
- Republish `0.9.2` without functional changes

## [0.9.2]
- Add `--header` option for passing custom request headers

## [0.9.1]
- Initial release
