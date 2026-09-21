### Hexlet tests and linter status:
[![Actions Status](https://github.com/araslanov-e/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/araslanov-e/go-project-244/actions)
[![CI](https://github.com/araslanov-e/go-project-244/actions/workflows/ci.yml/badge.svg)](https://github.com/araslanov-e/go-project-244/actions/workflows/ci.yml)

## Описание

`gendiff` — CLI-утилита, которая сравнивает два конфигурационных файла и показывает разницу между ними.

## Использование

Поддерживаются форматы JSON (`.json`) и YAML (`.yml`, `.yaml`); формат определяется по расширению файла.
Файлы могут содержать вложенные структуры — они сравниваются рекурсивно.

```bash
make build
./bin/gendiff testdata/fixture/nested1.json testdata/fixture/nested2.json
./bin/gendiff --format stylish testdata/fixture/nested1.yml testdata/fixture/nested2.yaml
```

Формат вывода задаётся флагом `--format` (`-f`). По умолчанию используется `stylish`:
`+` — добавленное значение, `-` — удалённое, без знака — значение не изменилось.

![gendiff demo](docs/gendiff.gif)
