import pytest

from backend.model import Label


data = {
  "id": 123,
  "name": "Alice",
  "isActive": True,
  "roles": ["admin", "editor", "user"],
  "profile": {
    "age": 30,
    "email": "alice@example.com",
    "address": {
      "street": "123 Main St",
      "city": "Springfield"
    }
  },
  "projects": [
    {
      "name": "Project X",
      "tasks": [
        {"title": "Design UI", "completed": True},
        {"title": "Implement API", "completed": False}
      ]
    },
    {
      "name": "Project Y",
      "tasks": []
    }
  ]
}

@pytest.mark.parametrize("key,op,val,result",[
    (".id", "==", "123", True),
    (".id", "==", "124", False),
    (".id", "!=", "123", False),
    (".id", "!=", "124", True),
    (".id", "like", r'1\d3', True),
    (".id", "like", r'2\d1', False),
    (".id", "not like", r'1\d3', False),
    (".id", "not like", r'2\d1', True),
    (".roles[]", "contains", 'admin', True),
    (".roles[]", "contains", 'test', False),
    (".roles[]", "not contains", 'admin', False),
    (".roles[]", "not contains", 'test', True),
])
def test_dupa(key, op, val, result):
    label = Label(key=key, op=op, val=val)

    assert label.matches(data) == result
