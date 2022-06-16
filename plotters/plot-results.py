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

throughput_files = natsorted([join(sys.argv[1], f) for f in listdir(sys.argv[1]) if isfile(join(sys.argv[1], f))])
pprint(throughput_files)

latency_files = natsorted([join(sys.argv[2], f) for f in listdir(sys.argv[1]) if isfile(join(sys.argv[2], f))])
pprint(latency_files)

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

  avg_throughput = throughput_series.quantile(0.9)
  latency_90th = latency_series.quantile(0.9) / 1e6

  file_desc = throuput_file.split('/')
  exp_desc = file_desc[len(file_desc) - 1][:-4].split('-')
  n_clients, total_threads = int(exp_desc[1]), int(exp_desc[0])
  threads_per_client = total_threads / n_clients

  result_data = result_data.append(DataFrame([[n_clients, threads_per_client, total_threads, avg_throughput, latency_90th]], columns=['n_clients', 'threads_per_client', 'total_threads', 'avg_throughput', 'latency_90th']), ignore_index=True)

# series = read_csv(
#   sys.argv[1],
#   sep=' ',
#   squeeze=True,
# )

# series['total_threads'] = series['client_nodes'] * series['threads_per_client']

# print(series)

plot_result_data = result_data.sort_values('avg_throughput')
plot_result_data.plot(x='avg_throughput', y='latency_90th')
pyplot.xlabel("Vazão (média)")
pyplot.ylabel("Latência (percentil 90%)")
head, tail = ntpath.split(sys.argv[1])
if not isdir(f"./csvs/summary/{head}"): makedirs(f"./csvs/summary/{head}")
pyplot.savefig(f"./csvs/summary/{head}/{tail}.png")
result_data.to_csv(f"./csvs/summary/{head}/{tail}.csv", header=True, sep=';')
# pyplot.show()
