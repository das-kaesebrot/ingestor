import logging
from os import listdir
from os.path import join
from ..constants.allowed_file_extensions import AllowedFileExtension
from ..constants.ingesting_mode import IngestingMode
from ..constants.heic_mode import HeicMode
from ..utils.heic import HeicConverter
from ..utils.filename import FilenameUtils


class Ingestor:
    _directory: str
    _output_directory: str
    _heic_mode: HeicMode
    
    _filename_utils: FilenameUtils

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
        
        self._filename_utils = FilenameUtils(date_pattern=date_pattern, keep_original_filename=keep_original_filename, person_suffix=person_suffix)

        self._logger = logging.getLogger(__name__)

    def _get_new_filenames(self) -> dict[str, str]:
        filenames = {}
        image_files = Ingestor._find_image_files_in_directory(self._directory, include_heic=(self._heic_mode == HeicMode.COPY), include_raw=False)
        video_files = Ingestor._find_video_files_in_directory(self._directory)
        
        
        for image_file in image_files:
            filenames[image_file] = join(self._output_directory, FilenameUtils.get_filename_for_image(image_file_path=image_file))
        
        for video_file in video_files:
            filenames[video_file] = join(self._output_directory, FilenameUtils.get_filename_for_video(video_file_path=video_file))
        
        return filenames
    
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
    def _find_video_files_in_directory(directory) -> list[str]:
        return [
            join(directory, f)
            for f in listdir(directory)
            if AllowedFileExtension.is_video(f)
        ]
