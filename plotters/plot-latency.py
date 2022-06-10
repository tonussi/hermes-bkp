#!/usr/bin/env python3
import sys

from pandas import read_csv, to_timedelta
from matplotlib import pyplot

series = read_csv(
  sys.argv[1],
  sep=' ',
  names=('time', 'latency'),
  parse_dates=['time'],
  squeeze=True,
  index_col=0
)
print(series.quantile(0.9) / 1e6)

series.plot()
pyplot.show()
