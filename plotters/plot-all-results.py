#!/usr/bin/env python3
import sys
from os import listdir
from os.path import isfile, join

from pandas import DataFrame, read_csv
from matplotlib import pyplot

scenarios = sys.argv[1:]

axes = ()

for sc in scenarios:
  throughput_path = join(sc, 'throughput')
  latency_path = join(sc, 'latency')

  throughput_files = [join(throughput_path, f) for f in listdir(throughput_path) if isfile(join(throughput_path, f))]
  latency_files = [join(latency_path, f) for f in listdir(latency_path) if isfile(join(latency_path, f))]

  print(throughput_files, latency_files)

  result_data = DataFrame(columns=['avg_throughput', 'latency_90th'])

  for (throuput_file, latency_file) in zip(throughput_files, latency_files):
    throughput_series = read_csv(
      throuput_file,
      sep=' ',
      names=('unix_timestamp', 'req/s'),
      squeeze=True,
      index_col=0
    )

    latency_series = read_csv(
      latency_file,
      sep=' ',
      names=('unix_timestamp', 'latency'),
      squeeze=True,
      index_col=0
    )

    avg_throughput = throughput_series.mean()
    # Convert ns to ms (10e6)
    latency_90th = latency_series.quantile(0.9) / 1e6

    result_data = result_data.append(DataFrame([[avg_throughput, latency_90th]], columns=['avg_throughput', 'latency_90th']), ignore_index=True)

  result_data = result_data.sort_values('avg_throughput')

  print(result_data)

  axes = (*axes, result_data['avg_throughput'], result_data['latency_90th'])

pyplot.ylim()
pyplot.xlabel("throughput (avg)")
pyplot.ylabel("latency (90th)")
pyplot.plot(*axes)
pyplot.show()
