#!/usr/bin/env python3
import ntpath
import sys
import warnings
from os import listdir, makedirs
from os.path import isfile, join
from pprint import pprint

from genericpath import isdir
from matplotlib import pyplot
from natsort import natsorted
from pandas import DataFrame, read_csv

warnings.simplefilter(action='ignore', category=FutureWarning)

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

    latency_90th = latency_series.mean() / 1e6

    file_desc = throuput_file.split('/')
    exp_desc = file_desc[len(file_desc) - 1][:-4].split('-')
    n_clients, total_threads = int(exp_desc[1]), int(exp_desc[0])
    threads_per_client = total_threads / n_clients

    result_data = result_data.append(DataFrame([[n_clients, threads_per_client, total_threads, avg_throughput, latency_90th]], columns=['n_clients', 'threads_per_client', 'total_threads', 'avg_throughput', 'latency_90th']), ignore_index=True)

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
