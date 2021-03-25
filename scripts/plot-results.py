#!/usr/bin/env python3
import sys
from os import listdir
from os.path import isfile, join

from pandas import DataFrame, read_csv, concat
from matplotlib import pyplot

throughput_files = [join(sys.argv[1], f) for f in listdir(sys.argv[1]) if isfile(join(sys.argv[1], f))]
latency_files = [join(sys.argv[2], f) for f in listdir(sys.argv[1]) if isfile(join(sys.argv[2], f))]

print(throughput_files, latency_files)

result_data = DataFrame(columns=['avg_throughput', 'latency_90th'])

for (throuput_file, latency_file) in zip(throughput_files, latency_files):
  throughput_series = read_csv(
    throuput_file,
    sep=' ',
    names=('date', 'time', 'req/s'),
    parse_dates=[['date', 'time']],
    squeeze=True,
    index_col=0
  )

  latency_series = read_csv(
    latency_file,
    sep=' ',
    names=('date', 'time', 'latency'),
    parse_dates=[['date', 'time']],
    squeeze=True,
    index_col=0
  )

  avg_throughput = throughput_series.mean()
  latency_90th = latency_series.quantile(0.9) / 1000

  result_data = result_data.append(DataFrame([[avg_throughput, latency_90th]], columns=['avg_throughput', 'latency_90th']), ignore_index=True)

print(result_data.sort_values('avg_throughput'))

result_data = result_data.sort_values('avg_throughput')

# series = read_csv(
#   sys.argv[1],
#   sep=' ',
#   squeeze=True,
# )

# series['total_threads'] = series['client_nodes'] * series['threads_per_client']

# print(series)

result_data.plot(x='avg_throughput', y='latency_90th')
pyplot.show()