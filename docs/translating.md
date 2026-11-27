# Translating

WordList embeds one shared set of UTF-8 locale files in both the server and CLI. English (`en.json`) is the source for keys and meaning; keys are never translated. The required locale set is English, Spanish, French, German, Italian, and Arabic. All six files live in `src/common/i18n/locales/` and are included in every binary.

## Contributing a translation

For an existing language, update only its locale JSON file. For a new language:

1. Copy `en.json` to `src/common/i18n/locales/{code}.json` and translate every value.
2. Fill in its `meta` values: language code, English name, native name, text direction, and locale version.
3. Add the language to the registry in `src/common/i18n/i18n.go`, including its direction and CLDR plural categories. Mark it required only if it is part of the minimum set.
4. Add a row to the language table in AI.md through the project's specification maintenance process; do not edit AI.md as part of a translation contribution.
5. Update the language count and list in the README and this guide.
6. Run the repository's required tests in its approved container environment.

The embedded locale files and registry must contain exactly the same language codes. Every locale must have the same ordinary message keys and plural message keys as English. Each plural message must contain every CLDR category registered for that language and may additionally contain `zero`; unsupported categories are rejected. `{count}` is optional in each category, while all other placeholders must match across categories and locales. `TranslatePlural` uses `golang.org/x/text` cardinal rules and honors an explicit `zero` value. `TranslateFormat` performs literal `{name}` substitution and leaves placeholders without values visible. `go test ./src/common/i18n/...` checks locale-file parity, metadata, empty values, category coverage, and interpolation placeholders; `make test` includes these checks.

## Translation conventions

- Preserve placeholders such as `{count}`, `{total}`, and `{datetime}` exactly. Their position may change to fit the language.
- Translate complete phrases. Keep consistent, neutral formality and avoid word-for-word constructions that sound unnatural.
- Follow the language's punctuation, quotation, and spacing conventions. French, for example, uses appropriate spacing before some punctuation marks.
- Use the language's registered text direction. Arabic is right-to-left; translated layouts must remain usable when mirrored.
- Keep product names, protocol names, code, filenames, configuration keys, URLs, and literal values such as `dark` and `auto` unchanged.
- Translation values should be valid, idiomatic text. Machine translation is useful when the translator can reliably review its output; otherwise request a fluent contributor.
- When English keys or values change, update every locale in the same change. Do not leave empty strings or temporary translation markers.

## Glossary

| English term | Usage |
|---|---|
| passphrase | The generated sequence of words; use the established language-specific security term |
| wordlist | The source collection of words; preserve the product name `WordList` |
| token | A secret credential used to authorize API requests |
| rotate | Replace a token or secret with a new one, invalidating the prior value |
| revoke | Disable a credential so it can no longer authorize requests |
| privacy | Personal-data handling and privacy rights |
| preferences | Per-browser display and language settings |

## Review

Check placeholder parity and key parity in the i18n tests. Review the translated interface at 320, 768, and 1280 pixels, including right-to-left direction where applicable. Confirm long text wraps, form labels remain associated with controls, and keyboard focus stays visible. Update the README locale count and list when a language is added or retired.

## Credits

There are no named human translator credits in the current locale metadata. Contributors may add their chosen name or handle to `meta.translators` and this page when they provide or review a translation.
