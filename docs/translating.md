# Translating

WordList embeds one shared set of UTF-8 locale files in both the server and CLI. English (`en.json`) is the single source for keys and meaning; translators never edit it. The current minimum set is English, Spanish, French, German, Italian, and Arabic (six locales). Files live in `src/common/i18n/locales/` and are included in every binary. New languages are welcome; they are not held for native-speaker review when a reliable translation is available.

## Contributing a translation

For an existing language, update only its locale JSON file. A contributor who wants credit can add their name or handle to that locale's `meta.translators`; maintainers update this page and the README language list/count when merging. For a new language:

1. Copy `en.json` to `src/common/i18n/locales/{code}.json` and translate every value.
2. Fill in its `meta` values: language code, English name, native name, text direction, and locale version.
3. Add the language to the registry in `src/common/i18n/i18n.go`, including its direction and CLDR plural categories. Mark it required only if it is part of the minimum set.
4. Add a row to the language table in AI.md through the project's specification maintenance process; do not edit AI.md as part of a translation contribution.
5. Update the language count and list in the README and this guide; the count must match the locale JSON files.
6. Run `make test`; Go tools run in the project's approved `casjaysdev/go:latest` container and the gate also checks shell scripts, coverage, and the basic compiled-binary behavior.

The locale files and language registry must contain exactly the same language codes. Every locale must have the same ordinary and plural message keys as English. Each plural message must contain every CLDR category registered for that language and may additionally contain `zero`; unsupported categories are rejected. `{count}` is optional in each category, while all other placeholders must match across categories and locales. `TranslatePlural` uses `golang.org/x/text` cardinal rules and honors an explicit `zero` value. `TranslateFormat` performs literal `{name}` substitution and leaves placeholders without values visible.

Locale files use UTF-8 JSON with the same structure and key order as `en.json`, two-space indentation, no trailing commas, and no comments. Translate values only: keep keys unchanged. Each value is a complete phrase, not a fragment joined to another string. Use the same formal register throughout a locale and write idiomatically rather than translating word for word. If English text is ambiguous or incomplete, report it for maintainers instead of working around it in a translation.

## Translation conventions

- Preserve placeholders such as `{count}`, `{total}`, and `{datetime}` exactly. Their position may change to fit the language.
- Include every placeholder used by the English value with the exact same spelling and braces.
- Provide all CLDR plural forms required for the language (`zero`, `one`, `two`, `few`, `many`, `other` as applicable); use the CLDR plural chart to confirm a language's categories.
- Translate complete phrases. Keep consistent, neutral formality and avoid word-for-word constructions that sound unnatural.
- Follow the language's punctuation, quotation, and spacing conventions. French, for example, uses appropriate spacing before some punctuation marks.
- Use the language's registered text direction. Arabic is right-to-left; translated layouts must remain usable when mirrored.
- Keep product names, protocol names, code, filenames, configuration keys, URLs, and literal values such as `dark` and `auto` unchanged.
- Translation values should be valid, idiomatic text. Machine translation is useful when the translator can reliably review its output; otherwise request a fluent contributor.
- When English keys or values change, update every locale in the same change. Do not leave empty strings or temporary translation markers.
- The AI assistant may draft translations when it is fluent enough to produce reliable, natural text. Human corrections are welcome at any time and take precedence over later machine wording. Low-resource languages wait for a human contributor when reliable translation is not possible.
- Longer strings can take roughly 40% more space than English. Review the interface at 320, 768, and 1280 pixels; check RTL layouts mirrored, text wrapping, form labels, and visible keyboard focus.
- For scripts that may not be covered by the system font stack (such as CJK, Devanagari, or Thai), coordinate a suitable font-stack entry and verify that long unspaced text wraps.

## Glossary

| English term | Usage |
|---|---|
| passphrase | The generated sequence of words; use the established language-specific security term |
| wordlist | The source collection of words; preserve the product name `WordList` |
| token | A secret credential used to authorize API requests |
| paste | Insert copied text into a form or command input |
| rotate | Replace a token or secret with a new one, invalidating the prior value |
| revoke | Disable a credential so it can no longer authorize requests |
| privacy | Personal-data handling and privacy rights |
| preferences | Per-browser display and language settings |
| linked token | A separately revocable token associated with an owner token |
| bearer token | A credential sent in the HTTP `Authorization: Bearer` header |
| public key | The shareable half of a cryptographic keypair |
| private key | The secret half of a cryptographic keypair |

## Review

Reviewers check placeholders, plural categories, glossary consistency, layout, RTL mirroring, and the locale count in the README and this guide. Native-speaker review is welcome but is not a merge requirement. AI assistants should keep locale files synchronized whenever English strings change; no translation-status marker is needed. Retire a language only at its maintainer's request.

## Credits

There are no named human translator credits in the current locale metadata. Contributors may add their chosen name or handle to `meta.translators` and this page when they provide or review a translation.
