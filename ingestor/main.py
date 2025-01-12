#!/usr/bin/env python3

import argparse
import json
import logging
from sys import version_info


def main():
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
        default=".",
    )

    parser.add_argument(
        "-o",
        "--output-directory",
        help="Output directory to place renamed files in",
        type=str,
        required=True,
        default=None,
    )

    parser.add_argument(
        "--dry-run",
        help="Don't actually rename/copy any files",
        action="store_true",
        required=False,
        default=False,
    )

    parser.add_argument(
        "-s",
        "--silent",
        help="Suppress non-error output (sets logging level to ERROR)",
        action="store_true",
        required=False,
        default=False,
    )

    args = parser.parse_args()

    if args.silent:
        args.loglevel = logging.getLevelName(logging.ERROR).lower()

    logging.basicConfig(
        format="[%(asctime)s] [%(name)s] [%(levelname)s] %(message)s",
        level=args.loglevel.upper(),
    )

    logging.getLogger("main").debug(f"Args:\n{json.dumps(vars(args), indent=4)}")

    logger = logging.getLogger()

    if args.dry_run:
        logger.warning(f"Dry run active, skipping destructive operations")

    try:
        pass

    except KeyboardInterrupt as e:
        logger.warning("Interrupted by SIGINT")

    except Exception as e:
        logger.exception("Exception occured")
        exit(1)


if __name__ == "__main__":
    main()
