from ingestor.utils.time_offset_parser import TimeOffsetParser


HOURS = 23
MINUTES = 15
SECONDS = 20
TOTAL_SECONDS = (HOURS * 60 * 60) + (MINUTES * 60) + SECONDS


def test_parse_anonymous_offset_string():
    offset_string = f"{HOURS}:{MINUTES}:{SECONDS}"

    result = TimeOffsetParser.parse(offset_string)

    assert result.total_seconds() == TOTAL_SECONDS


def test_parse_positive_offset_string():

    offset_string = f"+{HOURS}:{MINUTES}:{SECONDS}"

    result = TimeOffsetParser.parse(offset_string)

    assert result.total_seconds() == TOTAL_SECONDS


def test_parse_negative_offset_string():
    offset_string = f"-{HOURS}:{MINUTES}:{SECONDS}"

    result = TimeOffsetParser.parse(offset_string)

    assert result.total_seconds() == -TOTAL_SECONDS
