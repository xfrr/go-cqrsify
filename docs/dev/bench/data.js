window.BENCHMARK_DATA = {
  "lastUpdate": 1755611705134,
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
          "id": "5ca39be39ddc0a2479d3a0f785e85ddb229438d3",
          "message": "Refactor GitHub Actions workflow to include additional paths for gh-pages deployment",
          "timestamp": "2024-10-06T13:04:55+02:00",
          "tree_id": "4b1d6ee9fd07f7acbc4b9330a66ce152f9115d56",
          "url": "https://github.com/xfrr/go-cqrsify/commit/5ca39be39ddc0a2479d3a0f785e85ddb229438d3"
        },
        "date": 1728212726890,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 305,
            "unit": "ns/op",
            "extra": "3370294 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 321.2,
            "unit": "ns/op",
            "extra": "4183113 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 331.3,
            "unit": "ns/op",
            "extra": "4340506 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 236.9,
            "unit": "ns/op",
            "extra": "5397751 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 210.3,
            "unit": "ns/op",
            "extra": "5734088 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 217,
            "unit": "ns/op",
            "extra": "4908260 times\n4 procs"
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
          "id": "d4d897ee34e01521ee8b5817d435d1d6c0363771",
          "message": "Refactor GitHub Actions workflow to include additional paths for gh-pages deployment",
          "timestamp": "2024-10-06T13:08:36+02:00",
          "tree_id": "d54e95a5a9fb062344569627b82879fa4eda6e2f",
          "url": "https://github.com/xfrr/go-cqrsify/commit/d4d897ee34e01521ee8b5817d435d1d6c0363771"
        },
        "date": 1728212947453,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 280.2,
            "unit": "ns/op",
            "extra": "3888747 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 304.7,
            "unit": "ns/op",
            "extra": "4453842 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 295.9,
            "unit": "ns/op",
            "extra": "3608416 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 204.3,
            "unit": "ns/op",
            "extra": "5764802 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 209.9,
            "unit": "ns/op",
            "extra": "5281879 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 218.1,
            "unit": "ns/op",
            "extra": "5422257 times\n4 procs"
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
          "id": "ae3203f0bb9c7918dc8482fb59465d63d104b05d",
          "message": "Add event WithTime func option",
          "timestamp": "2024-10-26T00:19:05+02:00",
          "tree_id": "2a21a29b301b215a37ced072045e05d1e4d01072",
          "url": "https://github.com/xfrr/go-cqrsify/commit/ae3203f0bb9c7918dc8482fb59465d63d104b05d"
        },
        "date": 1729894777738,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 325.5,
            "unit": "ns/op",
            "extra": "4442887 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 304.6,
            "unit": "ns/op",
            "extra": "4431025 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 336.8,
            "unit": "ns/op",
            "extra": "4253014 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 219.4,
            "unit": "ns/op",
            "extra": "4998321 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 206.5,
            "unit": "ns/op",
            "extra": "6078726 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 239.3,
            "unit": "ns/op",
            "extra": "4418137 times\n4 procs"
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
          "id": "c5a2a779b4f6b1d3e4d66c88a15b73c590141880",
          "message": "chore: refactor NextChange function signature to use generic type parameters",
          "timestamp": "2024-10-26T03:06:53+02:00",
          "tree_id": "71215f2712939a6a6833d0b2a92c9ac458b1cc81",
          "url": "https://github.com/xfrr/go-cqrsify/commit/c5a2a779b4f6b1d3e4d66c88a15b73c590141880"
        },
        "date": 1729904855466,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 322.6,
            "unit": "ns/op",
            "extra": "4469137 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 250.1,
            "unit": "ns/op",
            "extra": "4033246 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 283.7,
            "unit": "ns/op",
            "extra": "3567632 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 170.3,
            "unit": "ns/op",
            "extra": "6253437 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 212,
            "unit": "ns/op",
            "extra": "4961460 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 214.9,
            "unit": "ns/op",
            "extra": "5531733 times\n4 procs"
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
          "id": "092cac2ad0c236db7328556bc6fdac7fe9198827",
          "message": "fix: refactor Hydrate function to use generic type parameters for the base aggregate",
          "timestamp": "2024-10-27T02:45:06+02:00",
          "tree_id": "62eae63de8e35bc0c66b23b49f462818b449cc27",
          "url": "https://github.com/xfrr/go-cqrsify/commit/092cac2ad0c236db7328556bc6fdac7fe9198827"
        },
        "date": 1729989948401,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 347.7,
            "unit": "ns/op",
            "extra": "2986657 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 296.9,
            "unit": "ns/op",
            "extra": "4707060 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 274.1,
            "unit": "ns/op",
            "extra": "3803070 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 217.8,
            "unit": "ns/op",
            "extra": "6662610 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 209.1,
            "unit": "ns/op",
            "extra": "4983106 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 211.3,
            "unit": "ns/op",
            "extra": "6973680 times\n4 procs"
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
          "id": "aa1afde5acccd7cbb55ed23f26688226d119ebed",
          "message": "feat: enhance CQRS implementation with generic request and response types, add new error for invalid request response, and introduce Query interface",
          "timestamp": "2024-11-13T23:24:51+01:00",
          "tree_id": "106f9859cde83f1bd8a784279c0b1496bac30f31",
          "url": "https://github.com/xfrr/go-cqrsify/commit/aa1afde5acccd7cbb55ed23f26688226d119ebed"
        },
        "date": 1731536731394,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 301.4,
            "unit": "ns/op",
            "extra": "3493174 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 293.5,
            "unit": "ns/op",
            "extra": "3514911 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 302.3,
            "unit": "ns/op",
            "extra": "3561825 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 188.3,
            "unit": "ns/op",
            "extra": "6271497 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 211.2,
            "unit": "ns/op",
            "extra": "5423739 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 178.3,
            "unit": "ns/op",
            "extra": "6209508 times\n4 procs"
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
          "id": "69236f7b698da1ea4e7628930ca104c8a6c850d9",
          "message": "refactor: update function signatures to improve type safety and consistency",
          "timestamp": "2024-11-15T23:00:13+01:00",
          "tree_id": "2b0846e8a69b2c1a6e1b8ab3ffd5f6e44e762d6a",
          "url": "https://github.com/xfrr/go-cqrsify/commit/69236f7b698da1ea4e7628930ca104c8a6c850d9"
        },
        "date": 1731708054706,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 277,
            "unit": "ns/op",
            "extra": "4054386 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 256.1,
            "unit": "ns/op",
            "extra": "3993692 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 272.3,
            "unit": "ns/op",
            "extra": "3702553 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 200.2,
            "unit": "ns/op",
            "extra": "5820766 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 218.4,
            "unit": "ns/op",
            "extra": "5753188 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 247.4,
            "unit": "ns/op",
            "extra": "4326973 times\n4 procs"
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
          "id": "10eb6d3c542b76c8fe1fe1ea38838c8a70e6f8a7",
          "message": "fix: error handling and use dispatcher interface",
          "timestamp": "2024-11-30T02:50:22+01:00",
          "tree_id": "af7aebef5e0dceb670c9ad634acabf5cad884bbf",
          "url": "https://github.com/xfrr/go-cqrsify/commit/10eb6d3c542b76c8fe1fe1ea38838c8a70e6f8a7"
        },
        "date": 1732931469839,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 263.6,
            "unit": "ns/op",
            "extra": "3828976 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 271,
            "unit": "ns/op",
            "extra": "4052581 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 328,
            "unit": "ns/op",
            "extra": "4308492 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 226.9,
            "unit": "ns/op",
            "extra": "5307282 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 201,
            "unit": "ns/op",
            "extra": "5611041 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 192.6,
            "unit": "ns/op",
            "extra": "6501966 times\n4 procs"
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
          "id": "7fc3dfdbf23d02feb2176938d498418a829f1283",
          "message": "refactor: update InMemoryBus implementation to simplify publish logic and remove unused timeout features",
          "timestamp": "2025-01-09T00:46:13+01:00",
          "tree_id": "de43af8e6d1035317bca7d269b9c58e171650be4",
          "url": "https://github.com/xfrr/go-cqrsify/commit/7fc3dfdbf23d02feb2176938d498418a829f1283"
        },
        "date": 1736380006722,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 291,
            "unit": "ns/op",
            "extra": "3541600 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 260,
            "unit": "ns/op",
            "extra": "4108177 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 275.4,
            "unit": "ns/op",
            "extra": "3702710 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 211,
            "unit": "ns/op",
            "extra": "5255559 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 229.3,
            "unit": "ns/op",
            "extra": "5190248 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 231.3,
            "unit": "ns/op",
            "extra": "4370418 times\n4 procs"
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
          "id": "089093058f4076720c41051ef5a1646a660fcc3a",
          "message": "fix: event context any method returns incorrect type",
          "timestamp": "2025-01-12T02:29:51+01:00",
          "tree_id": "fa3c56971c2d0466e742eee121324b3fe65ee4f2",
          "url": "https://github.com/xfrr/go-cqrsify/commit/089093058f4076720c41051ef5a1646a660fcc3a"
        },
        "date": 1736645429736,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 278.8,
            "unit": "ns/op",
            "extra": "3890424 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 303.9,
            "unit": "ns/op",
            "extra": "4417063 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 290,
            "unit": "ns/op",
            "extra": "3619926 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 200.5,
            "unit": "ns/op",
            "extra": "5913205 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 201.3,
            "unit": "ns/op",
            "extra": "5793742 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 219.8,
            "unit": "ns/op",
            "extra": "4811006 times\n4 procs"
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
          "id": "c1b48c9cd9d162e968564ad45d3da30572ce67d0",
          "message": "fix: correct type casting in Any method of BaseContext",
          "timestamp": "2025-01-12T02:41:44+01:00",
          "tree_id": "64ca1755f44ac82ccfd79e7735c459b0d421409d",
          "url": "https://github.com/xfrr/go-cqrsify/commit/c1b48c9cd9d162e968564ad45d3da30572ce67d0"
        },
        "date": 1736646177676,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchString",
            "value": 289.4,
            "unit": "ns/op",
            "extra": "3884582 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchInt",
            "value": 266.2,
            "unit": "ns/op",
            "extra": "3862402 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStruct",
            "value": 262.3,
            "unit": "ns/op",
            "extra": "3944276 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchStringer",
            "value": 211.8,
            "unit": "ns/op",
            "extra": "5542093 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchGoStringer",
            "value": 200.5,
            "unit": "ns/op",
            "extra": "5758519 times\n4 procs"
          },
          {
            "name": "BenchmarkCommandDispatch/CommandDispatchCommand",
            "value": 196.5,
            "unit": "ns/op",
            "extra": "5160570 times\n4 procs"
          }
        ]
      }
    ],
    "Benchmark Tests": [
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
          "id": "b655e42251842288a070f561ff740e4b55ffe52d",
          "message": "fix: modify benchmark tests directory",
          "timestamp": "2025-08-15T02:08:52+02:00",
          "tree_id": "b94678dafd5eaa3517c61b0e4a9d265269914ef5",
          "url": "https://github.com/xfrr/go-cqrsify/commit/b655e42251842288a070f561ff740e4b55ffe52d"
        },
        "date": 1755216576767,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkBus_Publish/buffer-size-1",
            "value": 587,
            "unit": "ns/op",
            "extra": "2006708 times\n4 procs"
          },
          {
            "name": "BenchmarkBus_Publish/buffer-size-10",
            "value": 910.5,
            "unit": "ns/op",
            "extra": "2093791 times\n4 procs"
          },
          {
            "name": "BenchmarkBus_Publish/buffer-size-100",
            "value": 517.9,
            "unit": "ns/op",
            "extra": "2312646 times\n4 procs"
          },
          {
            "name": "BenchmarkBus_Publish/buffer-size-500",
            "value": 508.3,
            "unit": "ns/op",
            "extra": "2509099 times\n4 procs"
          },
          {
            "name": "BenchmarkBus_Publish/buffer-size-1000",
            "value": 500.3,
            "unit": "ns/op",
            "extra": "2434384 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 80.95,
            "unit": "ns/op",
            "extra": "14921359 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 94.95,
            "unit": "ns/op",
            "extra": "13442749 times\n4 procs"
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
          "id": "ab3a79ccd580a8448c20037db8a185937c4c1c9a",
          "message": "chore: update Go version to 1.24.x in CI workflow",
          "timestamp": "2025-08-15T02:09:47+02:00",
          "tree_id": "f1ff3d131b6c449cf0e1634d3c1f9854f2ba1862",
          "url": "https://github.com/xfrr/go-cqrsify/commit/ab3a79ccd580a8448c20037db8a185937c4c1c9a"
        },
        "date": 1755216633462,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkBus_Publish/buffer-size-1",
            "value": 622.3,
            "unit": "ns/op",
            "extra": "1993161 times\n4 procs"
          },
          {
            "name": "BenchmarkBus_Publish/buffer-size-10",
            "value": 609.5,
            "unit": "ns/op",
            "extra": "2019673 times\n4 procs"
          },
          {
            "name": "BenchmarkBus_Publish/buffer-size-100",
            "value": 531.8,
            "unit": "ns/op",
            "extra": "2264372 times\n4 procs"
          },
          {
            "name": "BenchmarkBus_Publish/buffer-size-500",
            "value": 495.8,
            "unit": "ns/op",
            "extra": "2421854 times\n4 procs"
          },
          {
            "name": "BenchmarkBus_Publish/buffer-size-1000",
            "value": 495.8,
            "unit": "ns/op",
            "extra": "2424296 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 82.83,
            "unit": "ns/op",
            "extra": "14509312 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 86.48,
            "unit": "ns/op",
            "extra": "13652804 times\n4 procs"
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
          "id": "3cd1dfaf4d91d7307c1199e0bf5c30564351e778",
          "message": "refactor: aggregate and event sourcing implementation",
          "timestamp": "2025-08-19T14:29:51+02:00",
          "tree_id": "73ffb0e1e45efe9481834d14f2c4c8539723917d",
          "url": "https://github.com/xfrr/go-cqrsify/commit/3cd1dfaf4d91d7307c1199e0bf5c30564351e778"
        },
        "date": 1755606955327,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 83.48,
            "unit": "ns/op",
            "extra": "14497784 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 87.65,
            "unit": "ns/op",
            "extra": "13747004 times\n4 procs"
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
          "id": "5983411dd7f92cabc9fefd3703ee1445f514c9cd",
          "message": "refactor: aggregate and event sourcing implementation",
          "timestamp": "2025-08-19T15:54:37+02:00",
          "tree_id": "0e4a5dbec67b2e6df7033c67b34e61e5ad97fd17",
          "url": "https://github.com/xfrr/go-cqrsify/commit/5983411dd7f92cabc9fefd3703ee1445f514c9cd"
        },
        "date": 1755611704468,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 87.75,
            "unit": "ns/op",
            "extra": "14587977 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 85.75,
            "unit": "ns/op",
            "extra": "13781767 times\n4 procs"
          }
        ]
      }
    ]
  }
}