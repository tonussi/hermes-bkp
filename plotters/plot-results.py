#!/usr/bin/env python3
from genericpath import isdir
import ntpath
import sys
from os import listdir, makedirs
from os.path import isfile, join

from pandas import DataFrame, read_csv
from matplotlib import pyplot

throughput_files = [join(sys.argv[1], f) for f in listdir(sys.argv[1]) if isfile(join(sys.argv[1], f))]
latency_files = [join(sys.argv[2], f) for f in listdir(sys.argv[1]) if isfile(join(sys.argv[2], f))]

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
  latency_90th = latency_series.quantile(0.9) / 1e6

  file_desc = throuput_file.split('/')
  exp_desc = file_desc[len(file_desc) - 1][:-4].split('-')
  n_clients, total_threads = int(exp_desc[1]), int(exp_desc[0])
  threads_per_client = total_threads / n_clients

  result_data = result_data.append(DataFrame([[n_clients, threads_per_client, total_threads, avg_throughput, latency_90th]], columns=['n_clients', 'threads_per_client', 'total_threads', 'avg_throughput', 'latency_90th']), ignore_index=True)

# result_data = result_data.sort_values('avg_throughput')

print(result_data.to_csv())

# series = read_csv(
#   sys.argv[1],
#   sep=' ',
#   squeeze=True,
# )

# series['total_threads'] = series['client_nodes'] * series['threads_per_client']

# print(series)

result_data.plot(x='avg_throughput', y='latency_90th')
head, tail = ntpath.split(sys.argv[1])
if not isdir(f"./figs/summary/{head}"): makedirs(f"./figs/summary/{head}")
pyplot.savefig(f"./figs/summary/{head}/lat_vs_thr.png")
# pyplot.show()
