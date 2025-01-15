import datetime
import logging
from PIL import Image
from os.path import basename, splitext, getmtime


class FilenameUtils:
    # https://exiftool.org/TagNames/EXIF.html
    EXIF_TAG_ID_DATETIMEORIGINAL = 0x9003

    date_pattern: str
    keep_original_filename: bool
    person_suffix: str
    output_directory: str

    def __init__(
        self, date_pattern: str, keep_original_filename: bool, person_suffix: str
    ):
        self.date_pattern = date_pattern
        self.keep_original_filename = keep_original_filename
        self.person_suffix = person_suffix

    def get_filename_for_image(
        self,
        *,
        image_file_path: str,
    ):
        return FilenameUtils._get_filename_for_image(
            date_pattern=self.date_pattern,
            image_file_path=image_file_path,
            person_suffix=self.person_suffix,
            keep_original_filename=self.keep_original_filename,
        )

    def get_filename_for_video(
        self,
        *,
        video_file_path: str,
    ):
        return FilenameUtils._get_filename_for_video(
            date_pattern=self.date_pattern,
            video_file_path=video_file_path,
            person_suffix=self.person_suffix,
            keep_original_filename=self.keep_original_filename,
        )

    @staticmethod
    def _get_filename_for_image(
        *,
        date_pattern: str,
        image_file_path: str,
        person_suffix: str,
        keep_original_filename: bool = False,
    ):
        date = FilenameUtils._get_exif_date(image_file_path)
        
        if not date:
            logging.getLogger(__name__).warning(f"Couldn't get EXIF date from '{image_file_path}', using modification date instead")
            date = FilenameUtils._get_mtime(image_file_path)
        
        return FilenameUtils._get_formatted_filename(
            date=date,
            date_pattern=date_pattern,
            file_path=image_file_path,
            person_suffix=person_suffix,
            keep_original_filename=keep_original_filename,
        )

    @staticmethod
    def _get_filename_for_video(
        *,
        date_pattern: str,
        video_file_path: str,
        person_suffix: str,
        keep_original_filename: bool = False,
    ):
        date = FilenameUtils._get_mtime(file_path=video_file_path)
        return FilenameUtils._get_formatted_filename(
            date=date,
            date_pattern=date_pattern,
            file_path=video_file_path,
            person_suffix=person_suffix,
            keep_original_filename=keep_original_filename,
        )
    
    @staticmethod
    def _get_mtime(file_path: str) -> datetime.datetime:
        return datetime.datetime.fromtimestamp(getmtime(file_path))

    @staticmethod
    def _get_formatted_filename(
        *,
        date: datetime.datetime,
        date_pattern: str,
        file_path: str,
        person_suffix: str,
        keep_original_filename: bool = False,
    ):
        formatted_date = datetime.datetime.strftime(date, date_pattern)

        original_filename_suffix = (
            FilenameUtils.get_basename_without_extension(file_path)
            if keep_original_filename
            else None
        )

        filename = f"{formatted_date}_{person_suffix}{original_filename_suffix}"

        return filename

    @staticmethod
    def _get_exif_date(image_file_path: str) -> datetime.datetime:
        with Image.open(image_file_path) as image:
            exif = image.getexif()
            if not exif:
                raise Exception(f"File '{image_file_path}' does not have EXIF data.")

            return datetime.datetime.strptime(
                exif[FilenameUtils.EXIF_TAG_ID_DATETIMEORIGINAL], "%Y:%m:%d %H:%M:%S"
            )

    @staticmethod
    def get_basename_without_extension(file_path: str) -> str:
        return splitext(basename(file_path))[0]

    @staticmethod
    def get_file_extension(file_path: str) -> str:
        extension = splitext(file_path)[-1].lstrip(".")
        return extension
