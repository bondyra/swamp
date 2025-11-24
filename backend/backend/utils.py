import re

from backend.model import Label, Op


def get_matches(label: Label, allowed_values):
    if label.op == Op.EQUALS or label.op == Op.CONTAINS:
        return [v for v in allowed_values if v == label.val]
    if label.op == Op.NOT_EQUALS or label.op == Op.NOT_CONTAINS:
        return [v for v in allowed_values if v != label.val]
    r = re.compile(label.val)
    if label.op == Op.LIKE:
        return [v for v in allowed_values if r.match(v)]
    if label.op == Op.NOT_LIKE:
        return [v for v in allowed_values if not r.match(v)]
