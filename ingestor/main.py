#!/usr/bin/env python3

import logging
from datetime import timedelta
from .utils.ingestor import Ingestor
from .utils.time_offset_parser import TimeOffsetParser
from .constants.ingesting_mode import IngestingMode
from .constants.heic_mode import HeicMode

DEFAULT_DIRECTORY = "."
DEFAULT_OUTPUT_DIRECTORY = "Merged"
DEFAULT_DATE_PATTERN = r"%Y-%m-%d %H.%M.%S"
DEFAULT_DRY_RUN = False
DEFAULT_PERSON_SUFFIX = "J_H"
DEFAULT_KEEP_ORIGINAL_FILENAME = False
DEFAULT_MODE = IngestingMode.MOVE
DEFAULT_HEIC_MODE = HeicMode.CONVERT
DEFAULT_TIME_CORRECTION_OFFSET = timedelta(seconds=0)


def cli_entrypoint():
    import argparse
    import json
    from sys import version_info
    from time import perf_counter
    from datetime import timedelta

    start = perf_counter()

    # set up logging config via argparse
    # custom behaviour for python versions < 3.11 as the level names mapping func was only added to the logging lib in 3.11
    if version_info[1] >= 11:
        loglevel_mapping = logging.getLevelNamesMapping().keys()
    else:
        loglevel_mapping = logging._nameToLevel.keys()

    available_levels = [level.lower() for level in loglevel_mapping]
    available_levels.remove(logging.getLevelName(logging.NOTSET).lower())
    available_levels.remove(logging.getLevelName(logging.WARNING).lower())

    parser = argparse.ArgumentParser(
        description="Script for ingesting and renaming image and video files from vacations from different people",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    parser.add_argument(
        "-l",
        "--logging",
        help="Set the log level",
        dest="loglevel",
        type=str,
        choices=available_levels,
        default=logging.getLevelName(logging.INFO).lower(),
    )

    parser.add_argument(
        "-d",
        "--directory",
        help="Directory to ingest files from",
        type=str,
        required=False,
        default=DEFAULT_DIRECTORY,
    )

    parser.add_argument(
        "-o",
        "--output-directory",
        help="Output directory to place renamed files in",
        type=str,
        required=False,
        default=DEFAULT_OUTPUT_DIRECTORY,
    )

    parser.add_argument(
        "--dry-run",
        help="Don't actually rename/copy any files",
        action="store_true",
        required=False,
        default=DEFAULT_DRY_RUN,
    )

    parser.add_argument(
        "-s",
        "--silent",
        help="Suppress non-error output (sets logging level to ERROR)",
        action="store_true",
        required=False,
        default=False,
    )

    parser.add_argument(
        "--date-pattern",
        help="Date pattern for filenames",
        type=str,
        required=False,
        default=DEFAULT_DATE_PATTERN,
    )

    parser.add_argument(
        "--person-suffix",
        help="File suffix of the person the files are from. Example: Julian Handy -> J_H",
        type=str,
        required=False,
        default=DEFAULT_PERSON_SUFFIX,
    )

    parser.add_argument(
        "-k",
        "--keep-original-filename",
        help="Whether to keep the original filename as part of the new filename",
        action="store_true",
        required=False,
        default=DEFAULT_KEEP_ORIGINAL_FILENAME,
    )

    parser.add_argument(
        "-m",
        "--mode",
        help="Operation mode",
        type=IngestingMode,
        required=False,
        choices=IngestingMode.list(),
        default=DEFAULT_MODE,
    )

    parser.add_argument(
        "--heic-mode",
        help="HEIC operation mode. Whether to convert HEIC files to JPG before copying or copying as is.",
        type=HeicMode,
        required=False,
        choices=HeicMode.list(),
        default=DEFAULT_HEIC_MODE,
    )

    parser.add_argument(
        "--time-correction-offset",
        help="A correction offset in the format [+/-]HH:MM:SS to apply to the media files in the folder. Can either be positive (no prefix or +) or negative (-).",
        type=TimeOffsetParser.parse,
        required=False,
        metavar="+00:00:00",
        default=DEFAULT_TIME_CORRECTION_OFFSET,
    )

    args = parser.parse_args()

    if args.silent:
        args.loglevel = logging.getLevelName(logging.ERROR).lower()

    logging.basicConfig(
        format="[%(asctime)s] [%(name)s] [%(levelname)s] %(message)s",
        level=args.loglevel.upper(),
    )

    logging.getLogger(__name__).debug(
        f"Args:\n{json.dumps(vars(args), indent=4, default=str)}"
    )

    exit_code = ingest(**vars(args))

    end = perf_counter()
    logging.getLogger(__name__).info(f"Done in {timedelta(seconds=end-start)}")

    exit(exit_code)


def ingest(
    *,
    directory: str = DEFAULT_DIRECTORY,
    output_directory: str = DEFAULT_OUTPUT_DIRECTORY,
    dry_run: bool = DEFAULT_DRY_RUN,
    person_suffix: str = DEFAULT_PERSON_SUFFIX,
    keep_original_filename: str = DEFAULT_KEEP_ORIGINAL_FILENAME,
    date_pattern: str = DEFAULT_DATE_PATTERN,
    mode: IngestingMode = DEFAULT_MODE,
    heic_mode: HeicMode = DEFAULT_HEIC_MODE,
    time_correction_offset: timedelta = DEFAULT_TIME_CORRECTION_OFFSET,
    **kwargs,
) -> int | None:
    logger = logging.getLogger(__name__)

    if time_correction_offset != DEFAULT_TIME_CORRECTION_OFFSET:
        raise NotImplementedError("Time correction not implemented yet!")

    if dry_run:
        logger.warning(f"Dry run active, skipping destructive operations")

    try:
        ingestor = Ingestor(
            directory=directory,
            output_directory=output_directory,
            person_suffix=person_suffix,
            keep_original_filename=keep_original_filename,
            date_pattern=date_pattern,
            heic_mode=heic_mode,
            time_correction_offset=time_correction_offset,
        )

        ingestor.do_the_thing(mode=mode, dry_run=dry_run)

    except KeyboardInterrupt as e:
        logger.warning("Interrupted by SIGINT")

    except Exception as e:
        logger.exception("Exception occured")
        return 1


if __name__ == "__main__":
    cli_entrypoint()
