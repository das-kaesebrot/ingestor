import logging
from os import listdir
from os.path import join
from ..constants.allowed_file_extensions import AllowedFileExtension


class Ingestor:
    _directory: str
    _output_directory: str
    _person_suffix: str
    _keep_original_filename: bool
    _date_pattern: str

    _logger: logging.Logger

    def __init__(
        self,
        *,
        directory: str,
        output_directory: str,
        person_suffix: str,
        keep_original_filename: bool,
        date_pattern: str,
    ):
        self._directory = directory
        self._output_directory = output_directory
        self._person_suffix = person_suffix
        self._keep_original_filename = keep_original_filename
        self._date_pattern = date_pattern

        self._logger = logging.getLogger(__name__)

    def move_all(self):
        pass

    def copy_all(self):
        pass

    def get_new_filenames(self) -> dict[str, str]:
        pass

    @staticmethod
    def _get_image_files_in_directory(directory) -> list[str]:
        return [
            join(directory, f)
            for f in listdir(directory)
            if AllowedFileExtension.is_image(f)
        ]

    @staticmethod
    def _get_video_files_in_directory(directory) -> list[str]:
        return [
            join(directory, f)
            for f in listdir(directory)
            if AllowedFileExtension.is_video(f)
        ]
