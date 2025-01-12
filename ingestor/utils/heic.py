#!/usr/bin/env python3

import logging
from os import listdir, remove
from os.path import join, splitext, abspath, basename
from PIL import Image
from pillow_heif import register_heif_opener


class HeicConverter:
    _logger: logging.Logger

    input_directory: str
    output_directory: str

    def __init__(self, input_directory: str, output_directory: str | None = None):
        self.input_directory = abspath(input_directory)
        self.output_directory = (
            input_directory if not output_directory else output_directory
        )
        self._logger = logging.getLogger()

    def run_conversion(self, delete_source_files: bool = False):
        files = HeicConverter._find_heic_files(self.input_directory)

        self._logger.info(
            f"Found {len(files)} HEIC/HEIF files to convert in '{self.input_directory}'"
        )

        self._logger.info("Starting conversion")

        register_heif_opener()

        for file in files:
            jpg_file = join(self.output_directory, splitext(basename(file))[0] + ".JPG")

            self._logger.debug(f"Converting file '{file}' to JPEG file '{jpg_file}'")

            try:
                with Image.open(file) as image:
                    exif = image.getexif()
                    image.convert("RGB").save(jpg_file, exif=exif)

                    if delete_source_files:
                        self._logger.debug(f"Deleting '{file}'")
                        remove(file)

            except Exception as e:
                self._logger.exception(
                    f"Exception occured while converting '{file}' to JPEG"
                )

    @staticmethod
    def _find_heic_files(directory: str) -> list[str]:
        return [
            join(directory, f)
            for f in listdir(directory)
            if f.lower().endswith(".heic") or f.lower().endswith(".heif")
        ]
