window.BENCHMARK_DATA = {
  "lastUpdate": 1728212572329,
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
          "id": "4ceccc8a425ce6122430c388586559bb9f16a82e",
          "message": "Refactor GitHub Actions workflow to include additional paths for gh-pages deployment",
          "timestamp": "2024-10-06T06:06:22+02:00",
          "tree_id": "a1bf0ecdab62f7af4ecadbeb2359e8799574292a",
          "url": "https://github.com/xfrr/go-cqrsify/commit/4ceccc8a425ce6122430c388586559bb9f16a82e"
        },
        "date": 1728187612558,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 308.9,
            "unit": "ns/op",
            "extra": "3444291 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 286.2,
            "unit": "ns/op",
            "extra": "3757632 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 257.7,
            "unit": "ns/op",
            "extra": "4165728 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 200,
            "unit": "ns/op",
            "extra": "5986875 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 210.8,
            "unit": "ns/op",
            "extra": "5151770 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 231.3,
            "unit": "ns/op",
            "extra": "4860484 times\n4 procs"
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
          "id": "c34b50ae8e8914c76639ecd0f8cd1c162b9a4e17",
          "message": "Refactor GitHub Actions workflow to include additional paths for gh-pages deployment",
          "timestamp": "2024-10-06T13:02:18+02:00",
          "tree_id": "3b36724e3eb648d480ee14e34a1731937c517185",
          "url": "https://github.com/xfrr/go-cqrsify/commit/c34b50ae8e8914c76639ecd0f8cd1c162b9a4e17"
        },
        "date": 1728212571919,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 256.2,
            "unit": "ns/op",
            "extra": "4014649 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 300.8,
            "unit": "ns/op",
            "extra": "4538868 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 320.8,
            "unit": "ns/op",
            "extra": "4186312 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 205.3,
            "unit": "ns/op",
            "extra": "5407129 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 211.3,
            "unit": "ns/op",
            "extra": "5373889 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 206.1,
            "unit": "ns/op",
            "extra": "5768727 times\n4 procs"
          }
        ]
      }
    ]
  }
}