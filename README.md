# ingestor

This tool aims to create a workflow for ingesting files from different people to a centralized folder while keeping image attribution and creation dates as filename metadata.

## Usage

```
usage: __main__.py [-h] [-l {critical,fatal,error,warn,info,debug}] [-d DIRECTORY] [-o OUTPUT_DIRECTORY] [--dry-run] [-s]
                   [--date-pattern DATE_PATTERN] [--person-suffix PERSON_SUFFIX] [-k] [-m {move,copy}]
                   [--heic-mode {convert,copy}] [--time-correction-offset +00:00:00]

Script for ingesting and renaming image and video files from vacations from different people

options:
  -h, --help            show this help message and exit
  -l, --logging {critical,fatal,error,warn,info,debug}
                        Set the log level (default: info)
  -d, --directory DIRECTORY
                        Directory to ingest files from (default: .)
  -o, --output-directory OUTPUT_DIRECTORY
                        Output directory to place renamed files in (default: Merged)
  --dry-run             Don't actually rename/copy any files (default: False)
  -s, --silent          Suppress non-error output (sets logging level to ERROR) (default: False)
  --date-pattern DATE_PATTERN
                        Date pattern for filenames (default: %Y-%m-%d %H.%M.%S)
  --person-suffix PERSON_SUFFIX
                        File suffix of the person the files are from. Example: Julian Handy -> J_H (default: J_H)
  -k, --keep-original-filename
                        Whether to keep the original filename as part of the new filename (default: False)
  -m, --mode {move,copy}
                        Operation mode (default: move)
  --heic-mode {convert,copy}
                        HEIC operation mode. Whether to convert HEIC files to JPG before copying or copying as is. (default:
                        convert)
  --time-correction-offset +00:00:00
                        A correction offset to apply to the media files in the folder. Can either be positive (no prefix or +)
                        or negative (-). (default: 00:00:00)
```

# Open Source License Attribution

This application uses Open Source components. You can find the source code of their open source projects along with license information below. We acknowledge and are grateful to these developers for their contributions to open source.

### [Pillow](https://github.com/python-pillow/Pillow)
- Copyright (c) 2010 by Jeffrey A. Clark and contributors
- [MIT-CMU License](https://github.com/python-pillow/Pillow/blob/main/LICENSE)

### [pillow_heif](https://github.com/bigcat88/pillow_heif)
- Copyright (c) 2021-2023, Pillow-Heif contributors
- [BSD-3-Clause License](https://github.com/bigcat88/pillow_heif/blob/master/LICENSE.txt)

### [exif_py](https://github.com/ianare/exif-py)
- Copyright (c) 2002-2007 Gene Cash
- Copyright (c) 2007-2023 Ianaré Sévi and contributors
- [BSD-3-Clause License](https://github.com/ianare/exif-py/blob/develop/LICENSE.txt)

### [ffmpeg_python](https://github.com/kkroening/ffmpeg-python)
- Copyright (c) ffmpeg_python contributors
- [Apache-2.0 License](https://github.com/kkroening/ffmpeg-python/blob/master/LICENSE)

### [Flake8](https://github.com/PyCQA/flake8)
- Copyright (c) 2011-2013 Tarek Ziade <tarek@ziade.org>
- Copyright (c) 2012-2016 Ian Cordasco <graffatcolmingov@gmail.com>
- [Flake8 license (MIT)](https://github.com/PyCQA/flake8/blob/main/LICENSE)

### [pytest](https://github.com/pytest-dev/pytest)
- Copyright (c) 2004-present [Holger Krekel](https://github.com/hpk42) and [other contributors](https://github.com/pytest-dev/pytest/blob/main/AUTHORS)
- [MIT License](https://github.com/pytest-dev/pytest/blob/main/LICENSE)
