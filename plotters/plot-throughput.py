#!/usr/bin/env python3
import sys

from pandas import read_csv
from matplotlib import pyplot

series = read_csv(
  sys.argv[1],
  sep=' ',
  names=('time', 'req/s'),
  parse_dates=['time'],
  squeeze=True,
  index_col=0
)
print(series.mean())

series.plot()
pyplot.show()
