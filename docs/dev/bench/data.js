window.BENCHMARK_DATA = {
  "lastUpdate": 1759492812833,
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
          "id": "05c127996444f5d90511f06421d2f75400b486ac",
          "message": "feat: implement common value objects and validation logic",
          "timestamp": "2025-08-22T00:13:14+02:00",
          "tree_id": "cf8f905caf59209f6cde2986ced1a3307b29f9e1",
          "url": "https://github.com/xfrr/go-cqrsify/commit/05c127996444f5d90511f06421d2f75400b486ac"
        },
        "date": 1755814431179,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkEmailCreation",
            "value": 441.9,
            "unit": "ns/op",
            "extra": "2742781 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.75,
            "unit": "ns/op",
            "extra": "317081913 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.85,
            "unit": "ns/op",
            "extra": "29297953 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 84.03,
            "unit": "ns/op",
            "extra": "13966568 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 87.11,
            "unit": "ns/op",
            "extra": "13568209 times\n4 procs"
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
          "id": "1bfa70bc9220022782e5fcd001822df32e405ff1",
          "message": "refactor: move aggregate package to domain folder",
          "timestamp": "2025-08-22T00:16:32+02:00",
          "tree_id": "009985af8af4287b0090f14e8881de24209d9b5c",
          "url": "https://github.com/xfrr/go-cqrsify/commit/1bfa70bc9220022782e5fcd001822df32e405ff1"
        },
        "date": 1755814627830,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkEmailCreation",
            "value": 436.2,
            "unit": "ns/op",
            "extra": "2731146 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.734,
            "unit": "ns/op",
            "extra": "321224382 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.95,
            "unit": "ns/op",
            "extra": "30019281 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 82.45,
            "unit": "ns/op",
            "extra": "14476221 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 86.52,
            "unit": "ns/op",
            "extra": "13719334 times\n4 procs"
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
          "id": "01054f0ab891bc93a6ce38c3bc37004bff7fb0c0",
          "message": "feat: add Identifier value object",
          "timestamp": "2025-08-22T15:06:09+02:00",
          "tree_id": "3c57baaebb499dbb5a8b307c3c25bea1fd1bfeec",
          "url": "https://github.com/xfrr/go-cqrsify/commit/01054f0ab891bc93a6ce38c3bc37004bff7fb0c0"
        },
        "date": 1755868005493,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkEmailCreation",
            "value": 438.3,
            "unit": "ns/op",
            "extra": "2742350 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.738,
            "unit": "ns/op",
            "extra": "320442970 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 0.3119,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 0.3116,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 92.11,
            "unit": "ns/op",
            "extra": "13091850 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 3.735,
            "unit": "ns/op",
            "extra": "321111182 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.735,
            "unit": "ns/op",
            "extra": "320987728 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.96,
            "unit": "ns/op",
            "extra": "28994760 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 89.57,
            "unit": "ns/op",
            "extra": "11912856 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 87.29,
            "unit": "ns/op",
            "extra": "13571694 times\n4 procs"
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
          "id": "6aa28c82af474419ff88a78c281478f79c0f547f",
          "message": "feat: implement domain policy framework and add example",
          "timestamp": "2025-08-26T21:06:08+02:00",
          "tree_id": "ac3cc87f4add007989a5a45a21b249e2ca0d6905",
          "url": "https://github.com/xfrr/go-cqrsify/commit/6aa28c82af474419ff88a78c281478f79c0f547f"
        },
        "date": 1756235206675,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6235,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 73.26,
            "unit": "ns/op",
            "extra": "16429194 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 196.8,
            "unit": "ns/op",
            "extra": "6048156 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 446.9,
            "unit": "ns/op",
            "extra": "2683668 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.745,
            "unit": "ns/op",
            "extra": "320514764 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.39,
            "unit": "ns/op",
            "extra": "54831532 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.24,
            "unit": "ns/op",
            "extra": "57557265 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 88.09,
            "unit": "ns/op",
            "extra": "13725925 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.051,
            "unit": "ns/op",
            "extra": "296317963 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.073,
            "unit": "ns/op",
            "extra": "292240849 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.32,
            "unit": "ns/op",
            "extra": "29213604 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 82.25,
            "unit": "ns/op",
            "extra": "13851177 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 86.04,
            "unit": "ns/op",
            "extra": "13665345 times\n4 procs"
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
          "id": "443ccbd0e99571ba37c035c2c8ba2a2f83032a66",
          "message": "feat: add HasErrors method to MultiError for error presence checking",
          "timestamp": "2025-08-26T21:06:33+02:00",
          "tree_id": "579175e0378397424fbba72d37098aef307dcff5",
          "url": "https://github.com/xfrr/go-cqrsify/commit/443ccbd0e99571ba37c035c2c8ba2a2f83032a66"
        },
        "date": 1756235235080,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6239,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 73.19,
            "unit": "ns/op",
            "extra": "16475926 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 197.4,
            "unit": "ns/op",
            "extra": "6051604 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 438.9,
            "unit": "ns/op",
            "extra": "2732833 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.737,
            "unit": "ns/op",
            "extra": "320489706 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.25,
            "unit": "ns/op",
            "extra": "55518524 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.24,
            "unit": "ns/op",
            "extra": "55317736 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 88.01,
            "unit": "ns/op",
            "extra": "13589415 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.046,
            "unit": "ns/op",
            "extra": "296436709 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.058,
            "unit": "ns/op",
            "extra": "296229038 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.38,
            "unit": "ns/op",
            "extra": "29461993 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 83.28,
            "unit": "ns/op",
            "extra": "13732974 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 86.21,
            "unit": "ns/op",
            "extra": "13956942 times\n4 procs"
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
          "id": "c9dfa2dc7d001111ff2c915bf715ae6acd3f90c2",
          "message": "feat: add message envelope and encoding/decoding support",
          "timestamp": "2025-09-09T01:02:02+02:00",
          "tree_id": "3229ad19bb5402e368af6e46d7739aefa2ca23ab",
          "url": "https://github.com/xfrr/go-cqrsify/commit/c9dfa2dc7d001111ff2c915bf715ae6acd3f90c2"
        },
        "date": 1757372589488,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6231,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.87,
            "unit": "ns/op",
            "extra": "16429254 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 197.1,
            "unit": "ns/op",
            "extra": "6098199 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 438.2,
            "unit": "ns/op",
            "extra": "2740178 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.424,
            "unit": "ns/op",
            "extra": "349705726 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 21.81,
            "unit": "ns/op",
            "extra": "54330549 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.14,
            "unit": "ns/op",
            "extra": "55696878 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 87.17,
            "unit": "ns/op",
            "extra": "13755914 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.052,
            "unit": "ns/op",
            "extra": "295921449 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.054,
            "unit": "ns/op",
            "extra": "295841980 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.8,
            "unit": "ns/op",
            "extra": "29471558 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 80.35,
            "unit": "ns/op",
            "extra": "14731975 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 84.25,
            "unit": "ns/op",
            "extra": "14046028 times\n4 procs"
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
          "id": "2d7aadd4d9442d7850ba618b40c9142577bdf3fe",
          "message": "feat: update message bus to support topic-based dispatching and enhance event handling",
          "timestamp": "2025-09-09T20:15:10+02:00",
          "tree_id": "7242d3cb537c1d91a51857d50d71215c1a40f8b8",
          "url": "https://github.com/xfrr/go-cqrsify/commit/2d7aadd4d9442d7850ba618b40c9142577bdf3fe"
        },
        "date": 1757441753006,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6272,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.79,
            "unit": "ns/op",
            "extra": "16423302 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 197.1,
            "unit": "ns/op",
            "extra": "6094924 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 445.5,
            "unit": "ns/op",
            "extra": "2716513 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.424,
            "unit": "ns/op",
            "extra": "349998494 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.63,
            "unit": "ns/op",
            "extra": "54916888 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.24,
            "unit": "ns/op",
            "extra": "56020154 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 91.55,
            "unit": "ns/op",
            "extra": "13652084 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.1,
            "unit": "ns/op",
            "extra": "285634940 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.056,
            "unit": "ns/op",
            "extra": "295695405 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.46,
            "unit": "ns/op",
            "extra": "29621394 times\n4 procs"
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
          "id": "e6c15c194c0f4260428eb4b9d2c2df28d897c6b2",
          "message": "feat: update message bus to support topic-based dispatching and enhance event handling",
          "timestamp": "2025-09-09T20:21:25+02:00",
          "tree_id": "a7a7649e0f5839fab62e0e490d0946f118aadefc",
          "url": "https://github.com/xfrr/go-cqrsify/commit/e6c15c194c0f4260428eb4b9d2c2df28d897c6b2"
        },
        "date": 1757442122859,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6282,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.82,
            "unit": "ns/op",
            "extra": "16514451 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 197,
            "unit": "ns/op",
            "extra": "6107162 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 435.7,
            "unit": "ns/op",
            "extra": "2740496 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.427,
            "unit": "ns/op",
            "extra": "349966071 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.33,
            "unit": "ns/op",
            "extra": "59197608 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 21.1,
            "unit": "ns/op",
            "extra": "59385646 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 87.04,
            "unit": "ns/op",
            "extra": "13720779 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.049,
            "unit": "ns/op",
            "extra": "296167599 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.052,
            "unit": "ns/op",
            "extra": "295933358 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39,
            "unit": "ns/op",
            "extra": "29312062 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 54.84,
            "unit": "ns/op",
            "extra": "21905400 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 58.9,
            "unit": "ns/op",
            "extra": "20359581 times\n4 procs"
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
          "id": "5722726b9099875dd65ad2d37d719b6569a697a0",
          "message": "feat: update message bus to support topic-based dispatching and enhance event handling",
          "timestamp": "2025-09-09T20:24:54+02:00",
          "tree_id": "8e30f459f23c5241c097311e2857c70dda47a577",
          "url": "https://github.com/xfrr/go-cqrsify/commit/5722726b9099875dd65ad2d37d719b6569a697a0"
        },
        "date": 1757442333085,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6254,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.9,
            "unit": "ns/op",
            "extra": "16435186 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 197.4,
            "unit": "ns/op",
            "extra": "6038079 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 438.6,
            "unit": "ns/op",
            "extra": "2728034 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.429,
            "unit": "ns/op",
            "extra": "348412935 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.34,
            "unit": "ns/op",
            "extra": "55390005 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.26,
            "unit": "ns/op",
            "extra": "54828090 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 87.82,
            "unit": "ns/op",
            "extra": "13395672 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.054,
            "unit": "ns/op",
            "extra": "295916355 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.05,
            "unit": "ns/op",
            "extra": "295568214 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.16,
            "unit": "ns/op",
            "extra": "29344496 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 54.94,
            "unit": "ns/op",
            "extra": "21012696 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 57.98,
            "unit": "ns/op",
            "extra": "20453560 times\n4 procs"
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
          "id": "4426c19a337a928745dd4133cb1b4e839d054e6b",
          "message": "feat: add retry backoff package\n\n- Implemented a comprehensive retry backoff package that includes various backoff strategies (constant, exponential, and jitter).\n- Added support for result-aware retries and batch processing with shared budgets.\n- Introduced idempotency helpers to ensure at-least-once execution safety.\n- Included observability hooks for metrics and logging during retry attempts.\n- Created an in-memory deduplication store for idempotency tokens.\n- Enhanced error classification with flexible retry conditions.\n- Documented the package with a detailed README outlining features and usage.",
          "timestamp": "2025-09-10T22:40:10+02:00",
          "tree_id": "572b0c7061a5fc052b41e861b1814cd782e57754",
          "url": "https://github.com/xfrr/go-cqrsify/commit/4426c19a337a928745dd4133cb1b4e839d054e6b"
        },
        "date": 1757536855485,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6323,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.87,
            "unit": "ns/op",
            "extra": "16424056 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 202.2,
            "unit": "ns/op",
            "extra": "5736426 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 445.1,
            "unit": "ns/op",
            "extra": "2705776 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.432,
            "unit": "ns/op",
            "extra": "348720651 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.89,
            "unit": "ns/op",
            "extra": "53733207 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 21.04,
            "unit": "ns/op",
            "extra": "54855018 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 96.9,
            "unit": "ns/op",
            "extra": "13262277 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.078,
            "unit": "ns/op",
            "extra": "290636650 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.052,
            "unit": "ns/op",
            "extra": "295592840 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 40.8,
            "unit": "ns/op",
            "extra": "28926478 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 57.92,
            "unit": "ns/op",
            "extra": "21252099 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 69.09,
            "unit": "ns/op",
            "extra": "18011216 times\n4 procs"
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
          "id": "46aa089ea2bbf7813ba5abd612fab08454328ac5",
          "message": "feat: add retry backoff package\n\n- Implemented a comprehensive retry backoff package that includes various backoff strategies (constant, exponential, and jitter).\n- Added support for result-aware retries and batch processing with shared budgets.\n- Introduced idempotency helpers to ensure at-least-once execution safety.\n- Included observability hooks for metrics and logging during retry attempts.\n- Created an in-memory deduplication store for idempotency tokens.\n- Enhanced error classification with flexible retry conditions.\n- Documented the package with a detailed README outlining features and usage.",
          "timestamp": "2025-09-10T22:40:36+02:00",
          "tree_id": "50ebfd6d4f137d5e3f817d97f57b4f07eb45caf1",
          "url": "https://github.com/xfrr/go-cqrsify/commit/46aa089ea2bbf7813ba5abd612fab08454328ac5"
        },
        "date": 1757536870665,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6255,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.84,
            "unit": "ns/op",
            "extra": "16421850 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 200.7,
            "unit": "ns/op",
            "extra": "6019194 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 438.9,
            "unit": "ns/op",
            "extra": "2729222 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.425,
            "unit": "ns/op",
            "extra": "350068761 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.61,
            "unit": "ns/op",
            "extra": "54627232 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 21.35,
            "unit": "ns/op",
            "extra": "50517045 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 89.43,
            "unit": "ns/op",
            "extra": "13188357 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.056,
            "unit": "ns/op",
            "extra": "295311984 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.06,
            "unit": "ns/op",
            "extra": "294270398 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 40.54,
            "unit": "ns/op",
            "extra": "28495569 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 55.71,
            "unit": "ns/op",
            "extra": "21342493 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 59.05,
            "unit": "ns/op",
            "extra": "20339497 times\n4 procs"
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
          "id": "5dd8185bd093d608c438256f8873ae626dd88a51",
          "message": "docs: simplify usage documentation by removing detailed examples",
          "timestamp": "2025-09-11T08:51:41+02:00",
          "tree_id": "18c00fa6d425f5bb03d2dcb993ab37f576280e7c",
          "url": "https://github.com/xfrr/go-cqrsify/commit/5dd8185bd093d608c438256f8873ae626dd88a51"
        },
        "date": 1757573572104,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.624,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.78,
            "unit": "ns/op",
            "extra": "16493610 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 196.3,
            "unit": "ns/op",
            "extra": "6125379 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 440.6,
            "unit": "ns/op",
            "extra": "2756577 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.422,
            "unit": "ns/op",
            "extra": "350052980 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.47,
            "unit": "ns/op",
            "extra": "55950375 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.07,
            "unit": "ns/op",
            "extra": "56707448 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 87.05,
            "unit": "ns/op",
            "extra": "13487503 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.052,
            "unit": "ns/op",
            "extra": "295679524 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.056,
            "unit": "ns/op",
            "extra": "296093715 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.07,
            "unit": "ns/op",
            "extra": "29705030 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_Dispatch",
            "value": 54.76,
            "unit": "ns/op",
            "extra": "21501369 times\n4 procs"
          },
          {
            "name": "BenchmarkInMemoryBus_DispatchWithMiddleware",
            "value": 58.3,
            "unit": "ns/op",
            "extra": "20546431 times\n4 procs"
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
          "id": "34690e55599e291d1e8336c622c656f02822bacd",
          "message": "chore: remove unused dependency on github.com/stretchr/objx",
          "timestamp": "2025-09-17T01:38:18+02:00",
          "tree_id": "027bb1203dde9d5a7121ee5edc43d7edede38742",
          "url": "https://github.com/xfrr/go-cqrsify/commit/34690e55599e291d1e8336c622c656f02822bacd"
        },
        "date": 1758065965701,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6245,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.83,
            "unit": "ns/op",
            "extra": "16463541 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 196.3,
            "unit": "ns/op",
            "extra": "6131542 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 565.2,
            "unit": "ns/op",
            "extra": "2138658 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.427,
            "unit": "ns/op",
            "extra": "348783855 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.3,
            "unit": "ns/op",
            "extra": "55063686 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.2,
            "unit": "ns/op",
            "extra": "56735420 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 88.25,
            "unit": "ns/op",
            "extra": "13597046 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.14,
            "unit": "ns/op",
            "extra": "279898140 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 4.047,
            "unit": "ns/op",
            "extra": "296352762 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.33,
            "unit": "ns/op",
            "extra": "28669777 times\n4 procs"
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
          "id": "75cb5a4291764c76e9501ba5bd48840a642086db",
          "message": "feat: add JetstreamQueryBus implementation with request-response handling",
          "timestamp": "2025-09-17T19:43:43+02:00",
          "tree_id": "fc2991676fcc68040d6ca175272857278038ab66",
          "url": "https://github.com/xfrr/go-cqrsify/commit/75cb5a4291764c76e9501ba5bd48840a642086db"
        },
        "date": 1758131076731,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6413,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 73.23,
            "unit": "ns/op",
            "extra": "16441954 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 135.3,
            "unit": "ns/op",
            "extra": "8895027 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 444.7,
            "unit": "ns/op",
            "extra": "2671069 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.749,
            "unit": "ns/op",
            "extra": "320224243 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 21.8,
            "unit": "ns/op",
            "extra": "53375456 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.14,
            "unit": "ns/op",
            "extra": "57003621 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 89.74,
            "unit": "ns/op",
            "extra": "13298959 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.086,
            "unit": "ns/op",
            "extra": "293322102 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.734,
            "unit": "ns/op",
            "extra": "321166717 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.85,
            "unit": "ns/op",
            "extra": "30038610 times\n4 procs"
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
          "id": "c0445e72971de988ef840e448752f2451ca3a4dc",
          "message": "feat: update Go setup in CI workflow and add go.work configuration",
          "timestamp": "2025-09-17T20:53:00+02:00",
          "tree_id": "5fc57d6f60fbe8107e2979232cdc650124ada229",
          "url": "https://github.com/xfrr/go-cqrsify/commit/c0445e72971de988ef840e448752f2451ca3a4dc"
        },
        "date": 1758135234750,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6232,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.94,
            "unit": "ns/op",
            "extra": "16411212 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 134.2,
            "unit": "ns/op",
            "extra": "8962260 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 445.7,
            "unit": "ns/op",
            "extra": "2705258 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.741,
            "unit": "ns/op",
            "extra": "320262870 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.28,
            "unit": "ns/op",
            "extra": "53185194 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 22.85,
            "unit": "ns/op",
            "extra": "55739637 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 90.69,
            "unit": "ns/op",
            "extra": "13390182 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.046,
            "unit": "ns/op",
            "extra": "296727414 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.739,
            "unit": "ns/op",
            "extra": "320362232 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.02,
            "unit": "ns/op",
            "extra": "30122232 times\n4 procs"
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
          "id": "c1fbc399adbb1086c6fe5fcae19c3261c47e3b94",
          "message": "chore: update Codacy badge links",
          "timestamp": "2025-09-17T20:57:30+02:00",
          "tree_id": "f5a5121f2fc44cea13805d68c5c2410d76763bce",
          "url": "https://github.com/xfrr/go-cqrsify/commit/c1fbc399adbb1086c6fe5fcae19c3261c47e3b94"
        },
        "date": 1758135486471,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6224,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.97,
            "unit": "ns/op",
            "extra": "16451948 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 133.9,
            "unit": "ns/op",
            "extra": "8926892 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 447.6,
            "unit": "ns/op",
            "extra": "2678462 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.739,
            "unit": "ns/op",
            "extra": "320320915 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.3,
            "unit": "ns/op",
            "extra": "54376664 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.09,
            "unit": "ns/op",
            "extra": "56114074 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 88.98,
            "unit": "ns/op",
            "extra": "13474743 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.047,
            "unit": "ns/op",
            "extra": "296120844 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.742,
            "unit": "ns/op",
            "extra": "319350034 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.61,
            "unit": "ns/op",
            "extra": "30101618 times\n4 procs"
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
          "id": "56fab3d2cbd1265377c1d5550fbc09f2e49f9cac",
          "message": "ci: remove codacy scan schedule",
          "timestamp": "2025-09-17T21:07:36+02:00",
          "tree_id": "8e655b666be42417fdcf376b01204690f3f257c1",
          "url": "https://github.com/xfrr/go-cqrsify/commit/56fab3d2cbd1265377c1d5550fbc09f2e49f9cac"
        },
        "date": 1758136095656,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6242,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 73.05,
            "unit": "ns/op",
            "extra": "16433685 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 136.2,
            "unit": "ns/op",
            "extra": "8860365 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 444.2,
            "unit": "ns/op",
            "extra": "2431924 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.741,
            "unit": "ns/op",
            "extra": "319377246 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.42,
            "unit": "ns/op",
            "extra": "55914411 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.4,
            "unit": "ns/op",
            "extra": "54208329 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 94.32,
            "unit": "ns/op",
            "extra": "13298542 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.049,
            "unit": "ns/op",
            "extra": "295954004 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.745,
            "unit": "ns/op",
            "extra": "320625396 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 40.89,
            "unit": "ns/op",
            "extra": "29976279 times\n4 procs"
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
          "id": "923a550af18664008f160d4a04cec352fe4b5d57",
          "message": "feat: update Go setup in CI workflow to use Go 1.25.x",
          "timestamp": "2025-09-17T21:08:24+02:00",
          "tree_id": "f2d6613cd12a1da29b47374fb903b3aa73fcd2e0",
          "url": "https://github.com/xfrr/go-cqrsify/commit/923a550af18664008f160d4a04cec352fe4b5d57"
        },
        "date": 1758136140123,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6364,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.9,
            "unit": "ns/op",
            "extra": "16455871 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 133.8,
            "unit": "ns/op",
            "extra": "8929693 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 448.5,
            "unit": "ns/op",
            "extra": "2715615 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.747,
            "unit": "ns/op",
            "extra": "318632610 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.35,
            "unit": "ns/op",
            "extra": "54080650 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.14,
            "unit": "ns/op",
            "extra": "55080933 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 89.05,
            "unit": "ns/op",
            "extra": "13341548 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.051,
            "unit": "ns/op",
            "extra": "295481311 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.74,
            "unit": "ns/op",
            "extra": "320163066 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.74,
            "unit": "ns/op",
            "extra": "28268050 times\n4 procs"
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
          "id": "ce2d5b2974ec33bb72c325655ecb56d3549685a2",
          "message": "chore: add concurrency settings to Codacy security scan workflow",
          "timestamp": "2025-09-17T21:09:53+02:00",
          "tree_id": "516d87b933db414e202cc87e8f72adcd2c75c708",
          "url": "https://github.com/xfrr/go-cqrsify/commit/ce2d5b2974ec33bb72c325655ecb56d3549685a2"
        },
        "date": 1758136227536,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6232,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.87,
            "unit": "ns/op",
            "extra": "16477440 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 135.9,
            "unit": "ns/op",
            "extra": "8913618 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 445.1,
            "unit": "ns/op",
            "extra": "2537785 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.736,
            "unit": "ns/op",
            "extra": "320713473 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.6,
            "unit": "ns/op",
            "extra": "54140895 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.35,
            "unit": "ns/op",
            "extra": "54729386 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 90.43,
            "unit": "ns/op",
            "extra": "13540621 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.06,
            "unit": "ns/op",
            "extra": "294947218 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.738,
            "unit": "ns/op",
            "extra": "319785889 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.91,
            "unit": "ns/op",
            "extra": "29250742 times\n4 procs"
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
          "id": "f4b5374af4372ab94886c77f9be3becbfb4b0c5a",
          "message": "feat: enhance JetStream integration with updated stream configurations and context handling",
          "timestamp": "2025-09-17T21:40:40+02:00",
          "tree_id": "1da3183e8af1d5fd7f6811fe28b807eaea5b9084",
          "url": "https://github.com/xfrr/go-cqrsify/commit/f4b5374af4372ab94886c77f9be3becbfb4b0c5a"
        },
        "date": 1758138074837,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6488,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.93,
            "unit": "ns/op",
            "extra": "16431157 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 134.8,
            "unit": "ns/op",
            "extra": "8898952 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 443.1,
            "unit": "ns/op",
            "extra": "2707988 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.733,
            "unit": "ns/op",
            "extra": "321376268 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.27,
            "unit": "ns/op",
            "extra": "55097802 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.2,
            "unit": "ns/op",
            "extra": "55822362 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 94.98,
            "unit": "ns/op",
            "extra": "13372797 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.068,
            "unit": "ns/op",
            "extra": "291958860 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.74,
            "unit": "ns/op",
            "extra": "320871598 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.69,
            "unit": "ns/op",
            "extra": "30074632 times\n4 procs"
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
          "id": "195e77a83a14e1122a633a2d96b10e9e889130b5",
          "message": "feat: implement NATS JetStream command and event buses",
          "timestamp": "2025-09-17T23:15:01+02:00",
          "tree_id": "35f94a98e19ce705d57285c221fe6f818e8394ab",
          "url": "https://github.com/xfrr/go-cqrsify/commit/195e77a83a14e1122a633a2d96b10e9e889130b5"
        },
        "date": 1758143744379,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.624,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 73.08,
            "unit": "ns/op",
            "extra": "16513708 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 136.9,
            "unit": "ns/op",
            "extra": "8934145 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 441.6,
            "unit": "ns/op",
            "extra": "2727530 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.736,
            "unit": "ns/op",
            "extra": "320996553 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.35,
            "unit": "ns/op",
            "extra": "55528135 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.15,
            "unit": "ns/op",
            "extra": "55318124 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 93.61,
            "unit": "ns/op",
            "extra": "13227110 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.048,
            "unit": "ns/op",
            "extra": "295895440 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.738,
            "unit": "ns/op",
            "extra": "320329069 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 41.59,
            "unit": "ns/op",
            "extra": "29798463 times\n4 procs"
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
          "id": "5daa9078f494f1a9c360f7fe82ffb90105f525fb",
          "message": "fix: correct function signatures and naming in command and event bus implementations",
          "timestamp": "2025-09-17T23:35:30+02:00",
          "tree_id": "4f9a27d5af1ce661314716e1589d1b99f60568ed",
          "url": "https://github.com/xfrr/go-cqrsify/commit/5daa9078f494f1a9c360f7fe82ffb90105f525fb"
        },
        "date": 1758144974041,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6406,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.99,
            "unit": "ns/op",
            "extra": "16483833 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 135,
            "unit": "ns/op",
            "extra": "8901532 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 441.1,
            "unit": "ns/op",
            "extra": "2733625 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.759,
            "unit": "ns/op",
            "extra": "314578126 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.46,
            "unit": "ns/op",
            "extra": "55182079 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 21.64,
            "unit": "ns/op",
            "extra": "55288893 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 91.64,
            "unit": "ns/op",
            "extra": "12906250 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.055,
            "unit": "ns/op",
            "extra": "295023530 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.737,
            "unit": "ns/op",
            "extra": "320345594 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.77,
            "unit": "ns/op",
            "extra": "29646999 times\n4 procs"
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
          "id": "292057458641fa237883c4f5e7b13cca3cde5eba",
          "message": "refactor: update options parameter type to PubSubMessageBusOption in NATS message bus implementations",
          "timestamp": "2025-09-17T23:42:08+02:00",
          "tree_id": "01c67dff71f4e89532f715f488d12d9248e964ad",
          "url": "https://github.com/xfrr/go-cqrsify/commit/292057458641fa237883c4f5e7b13cca3cde5eba"
        },
        "date": 1758145363911,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6432,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 73.04,
            "unit": "ns/op",
            "extra": "16391563 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 135.3,
            "unit": "ns/op",
            "extra": "8958790 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 443,
            "unit": "ns/op",
            "extra": "2712080 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.742,
            "unit": "ns/op",
            "extra": "320844104 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.34,
            "unit": "ns/op",
            "extra": "54544442 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 21.09,
            "unit": "ns/op",
            "extra": "55433534 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 89.66,
            "unit": "ns/op",
            "extra": "13343799 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.076,
            "unit": "ns/op",
            "extra": "296195713 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.758,
            "unit": "ns/op",
            "extra": "315122457 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.46,
            "unit": "ns/op",
            "extra": "28897968 times\n4 procs"
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
          "id": "aa94709735c1e6a020e27455027e2f7a7bd10a8b",
          "message": "feat: add messaging unit tests and enhance buses unsubscribe behaviour",
          "timestamp": "2025-09-18T01:20:19+02:00",
          "tree_id": "84a0f26e83af886692ca0d2b1ad6226bc7392cbb",
          "url": "https://github.com/xfrr/go-cqrsify/commit/aa94709735c1e6a020e27455027e2f7a7bd10a8b"
        },
        "date": 1758151257791,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.625,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 73.03,
            "unit": "ns/op",
            "extra": "15964683 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 135.3,
            "unit": "ns/op",
            "extra": "8879412 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 447.9,
            "unit": "ns/op",
            "extra": "2705179 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.782,
            "unit": "ns/op",
            "extra": "319694266 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.37,
            "unit": "ns/op",
            "extra": "54034980 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.22,
            "unit": "ns/op",
            "extra": "56096650 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 89.64,
            "unit": "ns/op",
            "extra": "13438224 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.049,
            "unit": "ns/op",
            "extra": "296126676 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.926,
            "unit": "ns/op",
            "extra": "301119854 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 39.01,
            "unit": "ns/op",
            "extra": "29367691 times\n4 procs"
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
          "id": "059da01db388fbaddc26f254f43d693f3d812ed1",
          "message": "chore: separate bus interfaces and enhance examples",
          "timestamp": "2025-09-18T02:07:13+02:00",
          "tree_id": "af5a5365a216d5416fe778e2f12e48cb757595f1",
          "url": "https://github.com/xfrr/go-cqrsify/commit/059da01db388fbaddc26f254f43d693f3d812ed1"
        },
        "date": 1758154067906,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6225,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 73,
            "unit": "ns/op",
            "extra": "16473574 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 135.7,
            "unit": "ns/op",
            "extra": "8849808 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 442.9,
            "unit": "ns/op",
            "extra": "2680042 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.734,
            "unit": "ns/op",
            "extra": "321032655 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.31,
            "unit": "ns/op",
            "extra": "53930290 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.17,
            "unit": "ns/op",
            "extra": "55859030 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 89.11,
            "unit": "ns/op",
            "extra": "13309740 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.093,
            "unit": "ns/op",
            "extra": "288894348 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.738,
            "unit": "ns/op",
            "extra": "320580519 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.7,
            "unit": "ns/op",
            "extra": "29731552 times\n4 procs"
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
          "id": "2a518f41c7cd008e0ed9b2fb2404ca1e880442e6",
          "message": "fix: enhance reply subject generation in PublishRequest for better correlation of replies",
          "timestamp": "2025-09-18T19:06:33+02:00",
          "tree_id": "03f456ed9b19f0500f922bd004f40077f821847b",
          "url": "https://github.com/xfrr/go-cqrsify/commit/2a518f41c7cd008e0ed9b2fb2404ca1e880442e6"
        },
        "date": 1758215229419,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6306,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.99,
            "unit": "ns/op",
            "extra": "16449919 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 142.4,
            "unit": "ns/op",
            "extra": "8457651 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 444.7,
            "unit": "ns/op",
            "extra": "2707471 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.738,
            "unit": "ns/op",
            "extra": "320785812 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.51,
            "unit": "ns/op",
            "extra": "49623013 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.31,
            "unit": "ns/op",
            "extra": "53995723 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 89.68,
            "unit": "ns/op",
            "extra": "13204492 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.05,
            "unit": "ns/op",
            "extra": "295704349 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.818,
            "unit": "ns/op",
            "extra": "320691927 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 40.02,
            "unit": "ns/op",
            "extra": "29925628 times\n4 procs"
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
          "id": "24e4de1ac31aec46deb9217ca184b0484b67cbc3",
          "message": "feat: implement unit of work pattern with PostgreSQL support",
          "timestamp": "2025-09-25T19:33:58+02:00",
          "tree_id": "57678bab9bfeb7cd7590969150bcc368e6747063",
          "url": "https://github.com/xfrr/go-cqrsify/commit/24e4de1ac31aec46deb9217ca184b0484b67cbc3"
        },
        "date": 1758821694001,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6223,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.9,
            "unit": "ns/op",
            "extra": "16442866 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 136.6,
            "unit": "ns/op",
            "extra": "8949261 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 442.2,
            "unit": "ns/op",
            "extra": "2735252 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.747,
            "unit": "ns/op",
            "extra": "319929225 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.29,
            "unit": "ns/op",
            "extra": "55358382 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.18,
            "unit": "ns/op",
            "extra": "55584505 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 94.79,
            "unit": "ns/op",
            "extra": "13435824 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.069,
            "unit": "ns/op",
            "extra": "295941348 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.733,
            "unit": "ns/op",
            "extra": "321228032 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 42.37,
            "unit": "ns/op",
            "extra": "29324227 times\n4 procs"
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
          "id": "f1b62a10f22cc9b126dc3a4a9b2013d0a4dc848a",
          "message": "feat: add init script for PostgreSQL setup and refactor uow example",
          "timestamp": "2025-09-25T23:42:43+02:00",
          "tree_id": "f68ff0ebb313ad2dd4e05b755ea1006e90615f50",
          "url": "https://github.com/xfrr/go-cqrsify/commit/f1b62a10f22cc9b126dc3a4a9b2013d0a4dc848a"
        },
        "date": 1758836598735,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6225,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.9,
            "unit": "ns/op",
            "extra": "16460056 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 133.8,
            "unit": "ns/op",
            "extra": "8952660 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 439.9,
            "unit": "ns/op",
            "extra": "2720689 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.74,
            "unit": "ns/op",
            "extra": "320818698 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.41,
            "unit": "ns/op",
            "extra": "54221052 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.23,
            "unit": "ns/op",
            "extra": "55569670 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 88.91,
            "unit": "ns/op",
            "extra": "13382896 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.048,
            "unit": "ns/op",
            "extra": "295891936 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.738,
            "unit": "ns/op",
            "extra": "320675053 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.66,
            "unit": "ns/op",
            "extra": "30077246 times\n4 procs"
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
          "id": "953b614c24beb5f8de0b31b6ffe906033fb11905",
          "message": "feat: implement Coordinates value object",
          "timestamp": "2025-09-26T19:39:29+02:00",
          "tree_id": "224b1674e00efd73530ce554a8fbf7a48965d46c",
          "url": "https://github.com/xfrr/go-cqrsify/commit/953b614c24beb5f8de0b31b6ffe906033fb11905"
        },
        "date": 1758908414024,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6291,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.92,
            "unit": "ns/op",
            "extra": "16451788 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 136,
            "unit": "ns/op",
            "extra": "8659035 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailCreation",
            "value": 448.3,
            "unit": "ns/op",
            "extra": "2690032 times\n4 procs"
          },
          {
            "name": "BenchmarkEmailEquality",
            "value": 3.426,
            "unit": "ns/op",
            "extra": "349442647 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/string",
            "value": 20.8,
            "unit": "ns/op",
            "extra": "56061944 times\n4 procs"
          },
          {
            "name": "BenchmarkNewIdentifier/int",
            "value": 20.29,
            "unit": "ns/op",
            "extra": "58830148 times\n4 procs"
          },
          {
            "name": "BenchmarkString",
            "value": 96.01,
            "unit": "ns/op",
            "extra": "12311868 times\n4 procs"
          },
          {
            "name": "BenchmarkEquals",
            "value": 4.069,
            "unit": "ns/op",
            "extra": "291745846 times\n4 procs"
          },
          {
            "name": "BenchmarkValidate",
            "value": 3.747,
            "unit": "ns/op",
            "extra": "318276873 times\n4 procs"
          },
          {
            "name": "BenchmarkMoneyAddition",
            "value": 38.8,
            "unit": "ns/op",
            "extra": "29760489 times\n4 procs"
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
          "id": "3478dacfc5eb2f2b7c76e29672f679af3cd48143",
          "message": "chore: update Money value object to use amountCents and currencyISO for improved clarity and consistency",
          "timestamp": "2025-09-26T20:16:29+02:00",
          "tree_id": "b0d34b2d8924213099e81be9b619c99898212fe5",
          "url": "https://github.com/xfrr/go-cqrsify/commit/3478dacfc5eb2f2b7c76e29672f679af3cd48143"
        },
        "date": 1758910627329,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6242,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 73.6,
            "unit": "ns/op",
            "extra": "16444123 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 133.9,
            "unit": "ns/op",
            "extra": "8879094 times\n4 procs"
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
          "id": "a58c0137d3d41227c11dbd0811dcff529b160424",
          "message": "feat: add Equals method to Coordinates value object for comparison functionality",
          "timestamp": "2025-09-26T20:17:46+02:00",
          "tree_id": "fc22caeac85530424d33f01adb5439496ff7afc9",
          "url": "https://github.com/xfrr/go-cqrsify/commit/a58c0137d3d41227c11dbd0811dcff529b160424"
        },
        "date": 1758910699823,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6229,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.91,
            "unit": "ns/op",
            "extra": "16475860 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 134.3,
            "unit": "ns/op",
            "extra": "8976948 times\n4 procs"
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
          "id": "cf08bc3b084944a61143f8c544c54b6f3b5767a5",
          "message": "refactor: change value object constructors to return value types instead of pointers",
          "timestamp": "2025-09-26T22:53:25+02:00",
          "tree_id": "75d80266d8a314e57812d9d5057a4038a226723c",
          "url": "https://github.com/xfrr/go-cqrsify/commit/cf08bc3b084944a61143f8c544c54b6f3b5767a5"
        },
        "date": 1758920038124,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6296,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.9,
            "unit": "ns/op",
            "extra": "16461327 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 134.9,
            "unit": "ns/op",
            "extra": "8887453 times\n4 procs"
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
          "id": "956e7fb9625904498c398cd1fe7942f4477c3737",
          "message": "feat: add IsValid and IsZero methods to Sex value object; update tests for Gender and Sex value objects",
          "timestamp": "2025-09-26T23:02:44+02:00",
          "tree_id": "b3e4bb91205157c5f44318ad8565d32da7f58a67",
          "url": "https://github.com/xfrr/go-cqrsify/commit/956e7fb9625904498c398cd1fe7942f4477c3737"
        },
        "date": 1758920599650,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.6234,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.96,
            "unit": "ns/op",
            "extra": "16430889 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 134.8,
            "unit": "ns/op",
            "extra": "8894749 times\n4 procs"
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
          "id": "fb9293635b65db21c1d59e96093cd5b7955a7659",
          "message": "chore: upgrade deps",
          "timestamp": "2025-09-26T23:36:20+02:00",
          "tree_id": "7d032e8eed2f38095a1031c32fd22b6af1fd961a",
          "url": "https://github.com/xfrr/go-cqrsify/commit/fb9293635b65db21c1d59e96093cd5b7955a7659"
        },
        "date": 1758922635425,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.9335,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.76,
            "unit": "ns/op",
            "extra": "16484007 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 136.2,
            "unit": "ns/op",
            "extra": "8910548 times\n4 procs"
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
          "id": "de2eb87f52d28d526b35b32002a059f4b161a77a",
          "message": "feat: expose NewEventAggregateReference constructor",
          "timestamp": "2025-09-27T00:01:08+02:00",
          "tree_id": "9affb1f2d303e0c80ad0ec0a2c987aab9fb2e4f5",
          "url": "https://github.com/xfrr/go-cqrsify/commit/de2eb87f52d28d526b35b32002a059f4b161a77a"
        },
        "date": 1758924095670,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.9365,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.78,
            "unit": "ns/op",
            "extra": "16453467 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 134.1,
            "unit": "ns/op",
            "extra": "8942322 times\n4 procs"
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
          "id": "192424378e0a49f13ac1b6f88537a442a5922b43",
          "message": "feat: add domain WithEventTimestamp function to set event timestamp",
          "timestamp": "2025-09-27T00:05:04+02:00",
          "tree_id": "8490b4e1a3dca615e8f16f7f2fe2ce9698ee589f",
          "url": "https://github.com/xfrr/go-cqrsify/commit/192424378e0a49f13ac1b6f88537a442a5922b43"
        },
        "date": 1758924330854,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.9343,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.74,
            "unit": "ns/op",
            "extra": "16517202 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 134.8,
            "unit": "ns/op",
            "extra": "8923116 times\n4 procs"
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
          "id": "0d9a9e1c145e52dcb00d8c0b9d1ebc35cd5c2464",
          "message": "chore: implement JetStreamMessagePublisherConfig",
          "timestamp": "2025-09-29T00:37:05+02:00",
          "tree_id": "cafe7649d9796a396f6dc712a8d7b286af02ca64",
          "url": "https://github.com/xfrr/go-cqrsify/commit/0d9a9e1c145e52dcb00d8c0b9d1ebc35cd5c2464"
        },
        "date": 1759161154021,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.9363,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.82,
            "unit": "ns/op",
            "extra": "16501328 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 135.2,
            "unit": "ns/op",
            "extra": "8864786 times\n4 procs"
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
          "id": "21b909992abb0deba50276c8d0a8dba65375d117",
          "message": "chore: extract validation error messages for latitude and longitude in Coordinates",
          "timestamp": "2025-09-29T20:31:10+02:00",
          "tree_id": "7702b87cbb622f166726a2067776e817e3d72ed1",
          "url": "https://github.com/xfrr/go-cqrsify/commit/21b909992abb0deba50276c8d0a8dba65375d117"
        },
        "date": 1759170700799,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.9418,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.7,
            "unit": "ns/op",
            "extra": "16505632 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 135.9,
            "unit": "ns/op",
            "extra": "8908687 times\n4 procs"
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
          "id": "208b4372948bc19acb7ff42142ad4623586f1d61",
          "message": "feat: implement apix package with JSON:API helpers and HTTP request validation",
          "timestamp": "2025-10-02T20:04:42+02:00",
          "tree_id": "bb7119bccd2118dad48c8a029ccc1d48c158b9fe",
          "url": "https://github.com/xfrr/go-cqrsify/commit/208b4372948bc19acb7ff42142ad4623586f1d61"
        },
        "date": 1759428319268,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.9777,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.94,
            "unit": "ns/op",
            "extra": "16458621 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 134.2,
            "unit": "ns/op",
            "extra": "8893846 times\n4 procs"
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
          "id": "392189a051e3592d9b8947edf2000370aa5956c7",
          "message": "feat: implement messaging HTTP server\n\n- Renamed baseMessage to BaseMessage for consistency in naming conventions.\n- Updated BaseCommand and BaseEvent to use BaseMessage.\n- Introduced HTTPMessageServer for handling incoming messages via HTTP.\n- Added JSON:API message decoding support in HTTPMessageServer.\n- Implemented message registration and decoding mechanisms.\n- Enhanced error handling and validation in HTTPMessageServer.\n- Added comprehensive tests for HTTPMessageServer and message handling.\n- Updated go.mod and go.sum to include new dependencies for HTTP handling.",
          "timestamp": "2025-10-02T20:08:20+02:00",
          "tree_id": "4ec0e6185312afb3d0d43f914ad23821daf13cc7",
          "url": "https://github.com/xfrr/go-cqrsify/commit/392189a051e3592d9b8947edf2000370aa5956c7"
        },
        "date": 1759428547094,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.9487,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.77,
            "unit": "ns/op",
            "extra": "16489815 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 135,
            "unit": "ns/op",
            "extra": "8866675 times\n4 procs"
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
          "id": "5891bd6967e598695a9685805bafa5bd497af1f2",
          "message": "feat: implement messaging HTTP server\n\n- Renamed baseMessage to BaseMessage for consistency in naming conventions.\n- Updated BaseCommand and BaseEvent to use BaseMessage.\n- Introduced HTTPMessageServer for handling incoming messages via HTTP.\n- Added JSON:API message decoding support in HTTPMessageServer.\n- Implemented message registration and decoding mechanisms.\n- Enhanced error handling and validation in HTTPMessageServer.\n- Added comprehensive tests for HTTPMessageServer and message handling.\n- Updated go.mod and go.sum to include new dependencies for HTTP handling.",
          "timestamp": "2025-10-02T20:09:37+02:00",
          "tree_id": "3653504a1ca4e30f0f0a8ef0ef43fafad51bc5a4",
          "url": "https://github.com/xfrr/go-cqrsify/commit/5891bd6967e598695a9685805bafa5bd497af1f2"
        },
        "date": 1759428619963,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.8664,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 52.59,
            "unit": "ns/op",
            "extra": "22794758 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 105,
            "unit": "ns/op",
            "extra": "11367630 times\n4 procs"
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
          "id": "a5af759a3074d45969ad8aa367e9bf60605c66d5",
          "message": "feat: HTTP command handling and server structure\n\n- Renamed CommandHTTPServer to CommandHandler for clarity and consistency.\n- Introduced CommandServer to handle HTTP commands using Gin framework.\n- Added CreateBaseCommandFromSingleDocument function to convert JSON:API single documents into BaseCommand.\n- Implemented server configuration options for timeouts in http_server_options.go.\n- Updated MessageHandler to improve HTTP message handling and error management.\n- Enhanced message decoding functions to support JSON:API format.\n- Added new dependencies in go.mod for improved functionality.\n- Cleaned up and organized code for better readability and maintainability.",
          "timestamp": "2025-10-03T13:59:42+02:00",
          "tree_id": "72589ae1c04b877d02663419f918343c6b4e35f2",
          "url": "https://github.com/xfrr/go-cqrsify/commit/a5af759a3074d45969ad8aa367e9bf60605c66d5"
        },
        "date": 1759492812189,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkPolicyEvaluation",
            "value": 0.9344,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompositePolicyEvaluation",
            "value": 72.75,
            "unit": "ns/op",
            "extra": "16494417 times\n4 procs"
          },
          {
            "name": "BenchmarkPolicyEngine",
            "value": 134.7,
            "unit": "ns/op",
            "extra": "8843644 times\n4 procs"
          }
        ]
      }
    ]
  }
}