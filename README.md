### Hexlet tests and linter status:
[![Actions Status](https://github.com/araslanov-e/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/araslanov-e/go-project-244/actions)
[![CI](https://github.com/araslanov-e/go-project-244/actions/workflows/ci.yml/badge.svg)](https://github.com/araslanov-e/go-project-244/actions/workflows/ci.yml)

## Описание

`gendiff` — CLI-утилита, которая сравнивает два конфигурационных файла и показывает разницу между ними.

## Использование

Поддерживаются форматы JSON (`.json`) и YAML (`.yml`, `.yaml`); формат определяется по расширению файла.

```bash
make build
./bin/gendiff testdata/fixture/file1.json testdata/fixture/file2.json
./bin/gendiff testdata/fixture/file1.yml testdata/fixture/file2.yaml
```

```
{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}
```
