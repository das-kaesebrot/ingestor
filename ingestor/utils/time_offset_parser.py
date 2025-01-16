import argparse
import re
from datetime import timedelta, datetime


class TimeOffsetParser:
    REGEX_TIMEOFFSET_FORMAT = r"^(\+?|-)(\d{2}):(\d{2}):(\d{2})$"
    
    def __init__(self):
        pass
    
    @staticmethod
    def parse(value: str) -> timedelta:
        value = value.strip()
        if not re.match(pattern=TimeOffsetParser.REGEX_TIMEOFFSET_FORMAT, string=value):
            raise argparse.ArgumentTypeError("Time offset is not in correct format")
        
        positive = True
        if value.startswith("-"): positive = False        
        value = value.lstrip("+").lstrip("-")
        
        parsed_date: datetime
        
        try:
            parsed_date = datetime.strptime(value, "%H:%M:%S")
        except Exception as e:
            raise argparse.ArgumentTypeError("Failed parsing time offset!").with_traceback(e.__traceback__)
        
        delta = timedelta(hours=parsed_date.hour, minutes=parsed_date.minute, seconds=parsed_date.second)
        
        if not positive:
            zero = timedelta(seconds=0)
            
            delta = zero - delta
        
        return delta