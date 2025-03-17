from datetime import timedelta
import logging
from os import listdir
from os.path import join, expanduser
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
        time_correction_offset: timedelta
    ):        
        self._directory = expanduser(directory.strip().rstrip("/").rstrip("\\"))
        self._output_directory = expanduser(output_directory.strip().rstrip("/").rstrip("\\"))
        self._heic_mode = heic_mode

        self._filename_utils = FilenameUtils(
            date_pattern=date_pattern,
            keep_original_filename=keep_original_filename,
            person_suffix=person_suffix,
            correction_offset=time_correction_offset,
        )

        self._logger = logging.getLogger(__name__)

    def do_the_thing(self, mode: IngestingMode, dry_run: bool = False):
        filenames = self._get_new_filenames()

        if dry_run:
            self._logger.warning("Dry run activated")

            mode_string = "moved" if mode == IngestingMode.MOVE else "copied"

            for old_name, new_name in filenames.items():
                self._logger.debug(
                    f"File would be {mode_string}: '{old_name}' -> '{new_name}'"
                )

            return

        if mode == IngestingMode.COPY:
            Ingestor.copy_all(filenames)
        elif mode == IngestingMode.MOVE:
            Ingestor.move_all(filenames)
        else:
            raise ValueError(f"Unsupported mode '{mode}'")

    @staticmethod
    def move_all(filenames: dict[str, str]):
        from shutil import move

        logger = logging.getLogger(__name__)

        for old_name, new_name in filenames.items():
            logger.debug(f"Moving: '{old_name}' -> '{new_name}'")
            move(old_name, new_name)

    @staticmethod
    def copy_all(filenames: dict[str, str]):
        from shutil import copy2

        logger = logging.getLogger(__name__)

        for old_name, new_name in filenames.items():
            logger.debug(f"Copying: '{old_name}' -> '{new_name}'")
            copy2(old_name, new_name)

    def _get_new_filenames(self) -> dict[str, str]:
        filenames = {}
        image_files = Ingestor._find_image_files_in_directory(
            self._directory,
            include_heic=(self._heic_mode == HeicMode.COPY),
            include_raw=False,
        )
        video_files = Ingestor._find_video_files_in_directory(self._directory)

        self._logger.info(
            f"Found {len(image_files)} image files in '{self._directory}'"
        )
        self._logger.info(
            f"Found {len(video_files)} video files in '{self._directory}'"
        )
        self._logger.debug(f"{image_files=}")
        self._logger.debug(f"{video_files=}")

        for image_file in image_files:
            filenames[image_file] = join(
                self._output_directory,
                self._filename_utils.get_filename_for_image(image_file_path=image_file),
            )

        for video_file in video_files:
            filenames[video_file] = join(
                self._output_directory,
                self._filename_utils.get_filename_for_video(video_file_path=video_file),
            )

        return filenames

    def _handle_heic(self):
        if not self._heic_mode == HeicMode.CONVERT:
            return

        converter = HeicConverter(self._directory)
        converter.run_conversion(delete_source_files=False)

    @staticmethod
    def _find_image_files_in_directory(
        directory, include_heic: bool = False, include_raw: bool = False
    ) -> list[str]:
        return [
            join(directory, f)
            for f in listdir(directory)
            if AllowedFileExtension.is_image(
                f, include_heic=include_heic, include_raw=include_raw
            )
        ]

    @staticmethod
    def _find_video_files_in_directory(directory) -> list[str]:
        return [
            join(directory, f)
            for f in listdir(directory)
            if AllowedFileExtension.is_video(f)
        ]
