from collections import namedtuple
from dataclasses import dataclass
from typing import Any, List


@dataclass
class Query:
    module_name: str
    module: Any
    thing_type: str
    props: List[str]
    quiet: bool


Result = namedtuple("Result", "id content blob")
