#!/usr/bin/env python3
from genericpath import isdir
import ntpath
from os import makedirs
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
head, tail = ntpath.split(sys.argv[1])
if not isdir(f"./csvs/{head}"): makedirs(f"./csvs/{head}")
pyplot.savefig(f"./csvs/{sys.argv[1]}.png")
# pyplot.show()
