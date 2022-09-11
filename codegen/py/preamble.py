#region preamble

import operator

to_str = str
to_bool = bool

add_i = add_f = add_c = concat = operator.add
sub_i = sub_f = sub_c = operator.sub
mul_i = mul_f = mul_c = operator.mul
mod_i = mod_f = mod_c = operator.mod
div_i = div_c = operator.floordiv
div_f = operator.truediv

eq = operator.eq
or_ = operator.or_
and_ = operator.and_

gt_i = gt_f = gt_c = operator.gt
lt_i = lt_f = lt_c = operator.lt

def char_at(s: str, i: int) -> int:
  while i > len(s):
    i -= len(s)
  while i < 0:
    i += len(s)
  return bytes(s, "utf8")[i]

def substr(s: str, a: int, b: int) -> str:
  while a > len(s):
    a -= len(s)
  while a < 0:
    a += len(s)
  while b > len(s):
    b -= len(s)
  while b < 0:
    b += len(s)
  return s[a:b]

def char_append(s: str, c: int) -> str:
  return s + bytes([c]).decode()

str_len = len
ctoi = int

#endregion preamble

