import pytest

from backend.model import Label, Op
from backend.utils import get_matches


@pytest.mark.parametrize("label,values,results",[
    (Label(key="x", op=Op.EQUALS, val="a"), ["a", "b"], ["a"]),
    (Label(key="x", op=Op.CONTAINS, val="a"), ["a", "b"], ["a"]),
    (Label(key="x", op=Op.CONTAINS, val="x"), ["a", "b"], []),
    (Label(key="x", op=Op.EQUALS, val="x"), ["a", "b"], []),
    (Label(key="x", op=Op.NOT_EQUALS, val="x"), ["a", "b"], ["a","b"]),
    (Label(key="x", op=Op.NOT_CONTAINS, val="x"), ["a", "b"], ["a","b"]),
    (Label(key="x", op=Op.LIKE, val="a|b"), ["a", "b"], ["a", "b"]),
    (Label(key="x", op=Op.LIKE, val="(abc)|(def)"), ["ab", "abc", "cdef", "def"], ["abc","def"]),
    (Label(key="x", op=Op.LIKE, val="(abc)|(def)"), ["ab", "abc", "cdef"], ["abc"]),
    (Label(key="x", op=Op.NOT_LIKE, val="(abc)|(def)"), ["ab", "abc", "cdef", "def"], ["ab","cdef"]),
])
def test_dupa(label, values, results):
    assert get_matches(label, values) == results
