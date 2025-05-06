import argparse
import logging
import json
from ..main import ingest, DEFAULT_DIRECTORY, DEFAULT_OUTPUT_DIRECTORY
from sys import version_info
from tkinter import *
from tkinter import ttk

def gui_entrypoint():
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
        description="GUI tool for ingesting and renaming image and video files from vacations from different people",
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

    args = parser.parse_args()
    
    logging.basicConfig(
        format="[%(asctime)s] [%(name)s] [%(levelname)s] %(message)s",
        level=args.loglevel.upper(),
    )

    logging.getLogger(__name__).debug(
        f"Args:\n{json.dumps(vars(args), indent=4, default=str)}"
    )
    
    logger = logging.getLogger(__name__)
    
    try:
        root = Tk()
        frm = ttk.Frame(root, padding=10)
        frm.grid()
        ttk.Label(frm, text="Hello World!").grid(column=0, row=0)
        ttk.Button(frm, text="Quit", command=root.destroy).grid(column=1, row=0)
        root.mainloop()

    except KeyboardInterrupt as e:
        logger.warning("Interrupted by SIGINT")

    except Exception as e:
        logger.exception("Exception occured")
        return 1


if __name__ == "__main__":
    gui_entrypoint()
