import logging
from os import listdir
from os.path import join
from ..constants.allowed_file_extensions import AllowedFileExtension
from ..constants.heic_mode import HeicMode
from ..utils.heic import HeicConverter


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
        heic_mode: HeicMode,
    ):
        self._directory = directory
        self._output_directory = output_directory
        self._heic_mode = heic_mode
        self._keep_original_filename = keep_original_filename
        self._date_pattern = date_pattern

        self._logger = logging.getLogger(__name__)

    def move_all(self):
        pass

    def copy_all(self):
        pass

    
    def _handle_heic(self):
        if not self._heic_mode == HeicMode.CONVERT:
            return
        
        converter = HeicConverter(self._directory)
        converter.run_conversion(delete_source_files=False)
            

    @staticmethod
    def _find_image_files_in_directory(directory, include_heic: bool = False, include_raw: bool = False) -> list[str]:
        return [
            join(directory, f)
            for f in listdir(directory)
            if AllowedFileExtension.is_image(f, include_heic=include_heic, include_raw=include_raw)
        ]

    @staticmethod
    def _get_video_files_in_directory(directory) -> list[str]:
        return [
            join(directory, f)
            for f in listdir(directory)
            if AllowedFileExtension.is_video(f)
        ]
