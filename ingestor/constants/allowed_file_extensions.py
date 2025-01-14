from ..utils.filename import FilenameUtils


class AllowedFileExtension:
    # apple image extensions
    HEIC = ["heic", "heif"]

    # image extensions
    _JPEG = ["jpeg", "jpg", "jpe", "jif", "jfif", "jfi"]
    _PNG = ["png"]
    _TIFF = ["tiff", "tif"]
    _RAW = ["cr2", "cr3", "arw", "raw", "dng"]

    # video extensions
    _VIDEO = [
        "3g2",
        "3gp",
        "asf",
        "asx",
        "avi",
        "flv",
        "mkv",
        "mov",
        "mp4",
        "mpg",
        "mpeg",
        "mts",
        "ogv",
        "rm",
        "swf",
        "vob",
        "wmv",
        "webm",
    ]

    def __init__(self):
        pass

    @staticmethod
    def image(include_heic: bool = False, include_raw: bool = False) -> list[str]:
        extensions = [
            *AllowedFileExtension._JPEG,
            *AllowedFileExtension._PNG,
            *AllowedFileExtension._TIFF,
        ]

        if include_heic:
            [*extensions, *AllowedFileExtension.HEIC]

        if include_raw:
            [*extensions, *AllowedFileExtension._RAW]

        return extensions

    @staticmethod
    def video() -> list[str]:
        return AllowedFileExtension._VIDEO

    @staticmethod
    def is_image(
        file_path: str, include_heic: bool = False, include_raw: bool = False
    ) -> bool:
        return FilenameUtils.get_file_extension(
            file_path
        ).lower() in AllowedFileExtension.image(
            include_heic=include_heic, include_raw=include_raw
        )

    @staticmethod
    def is_video(file_path: str) -> bool:
        return (
            FilenameUtils.get_file_extension(file_path).lower()
            in AllowedFileExtension.video()
        )

    @staticmethod
    def is_heic(file_path: str) -> bool:
        return (
            FilenameUtils.get_file_extension(file_path).lower()
            in AllowedFileExtension.HEIC
        )
