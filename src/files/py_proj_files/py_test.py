# `pip install pytest`
# `pytest xxx_test.py`


def add(a, b):
    return a + b


def test_add():
    assert add(3, 2) == 5
    assert add(-1, 1) == 0
    assert add(-1, 1) == 1
