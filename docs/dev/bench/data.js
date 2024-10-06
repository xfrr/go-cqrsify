window.BENCHMARK_DATA = {
  "lastUpdate": 1728187466651,
  "repoUrl": "https://github.com/xfrr/go-cqrsify",
  "entries": {
    "CQRS Benchmark": [
      {
        "commit": {
          "author": {
            "email": "francisco.romero.1994@gmail.com",
            "name": "xfrr",
            "username": "xfrr"
          },
          "committer": {
            "email": "francisco.romero.1994@gmail.com",
            "name": "xfrr",
            "username": "xfrr"
          },
          "distinct": true,
          "id": "79eb745c02251ee157978167140b14629832eea6",
          "message": "Add github action to publish gh-pages",
          "timestamp": "2024-10-06T04:36:47+02:00",
          "tree_id": "6e04a7774616b17d1be265b426ed01bd366dce39",
          "url": "https://github.com/xfrr/go-cqrsify/commit/79eb745c02251ee157978167140b14629832eea6"
        },
        "date": 1728182514712,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 259.9,
            "unit": "ns/op",
            "extra": "3893216 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 316.7,
            "unit": "ns/op",
            "extra": "4325365 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 260.3,
            "unit": "ns/op",
            "extra": "3965358 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 195.2,
            "unit": "ns/op",
            "extra": "5938376 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 212.2,
            "unit": "ns/op",
            "extra": "4985049 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 208.1,
            "unit": "ns/op",
            "extra": "5448093 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "francisco.romero.1994@gmail.com",
            "name": "xfrr",
            "username": "xfrr"
          },
          "committer": {
            "email": "francisco.romero.1994@gmail.com",
            "name": "xfrr",
            "username": "xfrr"
          },
          "distinct": true,
          "id": "7dcc6f86e463e1154ea40254748fc09c544f3306",
          "message": "Refactor Makefile to include a test target and rename cover-out target to test-cover-out",
          "timestamp": "2024-10-06T06:03:55+02:00",
          "tree_id": "1586e4f0b29a4a6264672c95eb7bcf50218e78c7",
          "url": "https://github.com/xfrr/go-cqrsify/commit/7dcc6f86e463e1154ea40254748fc09c544f3306"
        },
        "date": 1728187466237,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 278,
            "unit": "ns/op",
            "extra": "3696748 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 308.3,
            "unit": "ns/op",
            "extra": "4291216 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 284.6,
            "unit": "ns/op",
            "extra": "3545474 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 205.9,
            "unit": "ns/op",
            "extra": "5520265 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 222.1,
            "unit": "ns/op",
            "extra": "5769296 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 189.1,
            "unit": "ns/op",
            "extra": "6202003 times\n4 procs"
          }
        ]
      }
    ]
  }
}