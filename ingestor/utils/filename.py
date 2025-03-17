import datetime
import exifread
import logging
import ffmpeg
from os.path import basename, splitext, getmtime


class FilenameUtils:
    # https://exiftool.org/TagNames/EXIF.html
    EXIF_TAG_ID_DATETIMEORIGINAL = 0x9003
    EXIF_TAG_NAME_DATETIMEORIGINAL = "EXIF DateTimeOriginal"

    date_pattern: str
    keep_original_filename: bool
    person_suffix: str
    output_directory: str
    correction_offset: datetime.timedelta

    def __init__(
        self, date_pattern: str, keep_original_filename: bool, person_suffix: str, correction_offset: datetime.timedelta
    ):
        self.date_pattern = date_pattern
        self.keep_original_filename = keep_original_filename
        self.person_suffix = person_suffix
        self.correction_offset = correction_offset

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
            time_correction_offset=self.correction_offset,
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
            time_correction_offset=self.correction_offset,
        )

    @staticmethod
    def _get_filename_for_image(
        *,
        date_pattern: str,
        image_file_path: str,
        person_suffix: str,
        time_correction_offset: datetime.timedelta,
        keep_original_filename: bool = False,
    ):
        date = FilenameUtils._get_exif_date(image_file_path)
        
        if not date:
            logging.getLogger(__name__).warning(f"Couldn't get EXIF date from '{image_file_path}', using file modification date instead")
            date = FilenameUtils._get_mtime(image_file_path)
        
        return FilenameUtils._get_formatted_filename(
            date=date,
            date_pattern=date_pattern,
            file_path=image_file_path,
            person_suffix=person_suffix,
            keep_original_filename=keep_original_filename,
            time_correction_offset=time_correction_offset
        )

    @staticmethod
    def _get_filename_for_video(
        *,
        date_pattern: str,
        video_file_path: str,
        person_suffix: str,
        time_correction_offset: datetime.timedelta,
        keep_original_filename: bool = False,
    ):
        date = FilenameUtils._get_video_creation_date(video_file_path)
        
        if not date:
            logging.getLogger(__name__).warning(f"Couldn't get video creation date from '{video_file_path}', using file modification date instead")
            date = FilenameUtils._get_mtime(video_file_path)
        
        return FilenameUtils._get_formatted_filename(
            date=date,
            date_pattern=date_pattern,
            file_path=video_file_path,
            person_suffix=person_suffix,
            keep_original_filename=keep_original_filename,
            time_correction_offset=time_correction_offset
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
        time_correction_offset: datetime.timedelta,
        keep_original_filename: bool = False,
    ):
        date = date + time_correction_offset        
        formatted_date = datetime.datetime.strftime(date, date_pattern)

        original_filename_suffix = (
            FilenameUtils.get_basename_without_extension(file_path)
            if keep_original_filename
            else ''
        )
        
        extension = FilenameUtils.get_file_extension(file_path)

        filename = f"{formatted_date}_{person_suffix}{original_filename_suffix}.{extension}"

        return filename
    
    @staticmethod
    def _get_video_creation_date(video_file_path: str) -> datetime.datetime:
        try:
            probe_result = ffmpeg.probe(video_file_path)
            creation_time_str = probe_result.get("format").get("tags").get("creation_time")
            return datetime.datetime.fromisoformat(creation_time_str)
        except Exception as e:
            logging.getLogger(__name__).exception(f"Error while probing '{video_file_path}'")
        
    @staticmethod
    def _get_exif_date(image_file_path: str) -> datetime.datetime:
        with open(image_file_path, "rb") as file_handle:
            tags = exifread.process_file(file_handle, stop_tag='DateTimeOriginal')
            
            if FilenameUtils.EXIF_TAG_NAME_DATETIMEORIGINAL in tags.keys():
                value = tags[FilenameUtils.EXIF_TAG_NAME_DATETIMEORIGINAL].values
                
                return datetime.datetime.strptime(value, "%Y:%m:%d %H:%M:%S")
            
            logging.getLogger(__name__).warning(f"Couldn't find tag '{FilenameUtils.EXIF_TAG_NAME_DATETIMEORIGINAL}'")

    @staticmethod
    def get_basename_without_extension(file_path: str) -> str:
        return splitext(basename(file_path))[0]

    @staticmethod
    def get_file_extension(file_path: str) -> str:
        extension = splitext(file_path)[-1].lstrip(".")
        return extension
