#!/usr/bin/env python3
import ntpath
import sys
import warnings
from os import makedirs

from genericpath import isdir
from matplotlib import pyplot
from pandas import read_csv

warnings.simplefilter(action='ignore', category=FutureWarning)

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
head, tail = ntpath.split(sys.argv[1])
if not isdir(f"./csvs/{head}"): makedirs(f"./csvs/{head}")
pyplot.savefig(f"./csvs/{sys.argv[1]}.png")
# pyplot.show()
