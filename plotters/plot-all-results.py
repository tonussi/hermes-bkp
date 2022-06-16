#!/usr/bin/env python3
from genericpath import isdir
import ntpath
import sys
from os import listdir, makedirs
from os.path import isfile, join
from pandas import DataFrame, read_csv
from matplotlib import pyplot
from natsort import natsorted
from pprint import pprint

root_scenarios = sys.argv[1]
scenarios = natsorted([join(root_scenarios, d) for d in listdir(root_scenarios) if isdir(join(root_scenarios, d))])

axes = ()

for sc in scenarios:
  throughput_path = join(sc, 'throughput')
  latency_path = join(sc, 'latency')

  throughput_files = natsorted([join(throughput_path, f) for f in listdir(throughput_path) if isfile(join(throughput_path, f))])
  pprint(throughput_files)

  latency_files = natsorted([join(latency_path, f) for f in listdir(latency_path) if isfile(join(latency_path, f))])
  pprint(latency_files)

  result_data = DataFrame(columns=['avg_throughput', 'latency_90th'])
  axes = ()

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

    latency_90th = latency_series.quantile(0.9) / 1e6

    result_data = result_data.append(DataFrame([[avg_throughput, latency_90th]], columns=['avg_throughput', 'latency_90th']), ignore_index=True)

  # plot_result_data = result_data.sort_values('avg_throughput')
  axes = (*axes, result_data['avg_throughput'], result_data['latency_90th'])

  pyplot.ylim()
  pyplot.xlabel("Vazão (média)")
  pyplot.ylabel("Latência (percentil 90%)")
  # pyplot.xticks(numpy.arange(min(axes[0]), max(axes[0]), 10.0))
  # pyplot.yticks(numpy.arange(min(axes[0]), max(axes[0]), 10.0))
  pyplot.plot(*axes)
  head, tail = ntpath.split(throughput_path)
  if not isdir(f"./csvs/summary/{head}"): makedirs(f"./csvs/summary/{head}")
  pyplot.savefig(f"./csvs/summary/{head}/lat_vs_vaz.png")
  result_data.to_csv(f"./csvs/summary/{head}/lat_vs_vaz.csv", header=True, sep=';')
  pyplot.cla()
  pyplot.clf()
