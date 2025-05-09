from datetime import timedelta
from zoneinfo import ZoneInfo
from ..constants.ingesting_mode import IngestingMode
from ..constants.heic_mode import HeicMode

class IngestorDefaultSettings:
    DIRECTORY = "."
    OUTPUT_DIRECTORY = "Merged"
    DATE_PATTERN = r"%Y-%m-%d %H.%M.%S"
    DRY_RUN = False
    PERSON_SUFFIX = "J_H"
    KEEP_ORIGINAL_FILENAME = False
    MODE = IngestingMode.MOVE
    HEIC_MODE = HeicMode.CONVERT
    TIME_CORRECTION_OFFSET = timedelta(seconds=0)
    TIMEZONE = ZoneInfo("Europe/Berlin")
    