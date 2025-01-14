import datetime
from PIL import Image
from os.path import basename, splitext


class FilenameUtils:
    # https://exiftool.org/TagNames/EXIF.html
    EXIF_TAG_ID_DATETIMEORIGINAL = 0x9003

    def __init__(self):
        pass

    @staticmethod
    def get_filename_for_image(
        *,
        date_pattern: str,
        image_file_path: str,
        person_suffix: str,
        keep_original_filename: bool = False,
    ):
        date = FilenameUtils.get_exif_date(image_file_path)
        formatted_date = datetime.datetime.strftime(date, date_pattern)

        original_filename_suffix = (
            FilenameUtils.get_basename_without_extension(image_file_path)
            if keep_original_filename
            else None
        )

        filename = f"{formatted_date}_{person_suffix}{original_filename_suffix}"

        return filename

    @staticmethod
    def get_exif_date(image_file_path: str) -> datetime:
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
