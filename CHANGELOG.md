# Changelog

## 1.0.0 (2026-04-22)


### Features

* adds Dockerfile ([f63d2a6](https://github.com/steamban/one2n-sre-bootcamp/commit/f63d2a6d98bb7bd80c8fac3a77619352c9f34208))
* adds local dev setup ([36ad127](https://github.com/steamban/one2n-sre-bootcamp/commit/36ad127004b4d4c25781b41e5d2e4ef4f4868bfc))
* **api:** adds crud apis for student ([59d78d2](https://github.com/steamban/one2n-sre-bootcamp/commit/59d78d2529ad0454130d6e68f1cbe08ae7b15806))
* **ci:** add github actions ci workflow ([4ef47d1](https://github.com/steamban/one2n-sre-bootcamp/commit/4ef47d1389d7552c4022c37a1e2ba47b99a319fb))
* **ci:** add hadolint check ([e772c61](https://github.com/steamban/one2n-sre-bootcamp/commit/e772c61feb174e982aa4c7f0a6c5dbf211112261))
* **db:** add students table ([77ad289](https://github.com/steamban/one2n-sre-bootcamp/commit/77ad2890751110615ebea8ff02ae844fd5fee8a6))
* **db:** update db config and connection setup ([19a0a83](https://github.com/steamban/one2n-sre-bootcamp/commit/19a0a839cb910303cbb888a181b7c83dd77ed4c1))
* **db:** updates student table ([2e62c9f](https://github.com/steamban/one2n-sre-bootcamp/commit/2e62c9fb0fc6051c499133d86d1feaba191aacfa))
* **release:** add release event trigger for docker builds ([0fb99c1](https://github.com/steamban/one2n-sre-bootcamp/commit/0fb99c11777c0169be707432af4251517ba77aaf))
* **release:** decouple docker job from release-please and use release tag for image versioning ([f067ad2](https://github.com/steamban/one2n-sre-bootcamp/commit/f067ad2f4fc5eb8ee1afc99d6a40271a7841c705))
* **release:** restrict release workflow to main branch only ([2fa2e87](https://github.com/steamban/one2n-sre-bootcamp/commit/2fa2e8799dd09f61112d3bb1b820aa10af607493))
* setup migration tools and db setup ([f97486b](https://github.com/steamban/one2n-sre-bootcamp/commit/f97486b7bbcd6bab9bd0147aedff3c51edc29886))


### Bug Fixes

* (api) adds validation struct for update student details API ([b19b172](https://github.com/steamban/one2n-sre-bootcamp/commit/b19b1723ad3a190f7318479766d010e2790ec1d9))
* adds .dockerignore ([95b44c2](https://github.com/steamban/one2n-sre-bootcamp/commit/95b44c2fe5ecae30e36bad93db0733193cde1c9c))
* adds digest for postgres ([d65f432](https://github.com/steamban/one2n-sre-bootcamp/commit/d65f432889d543d66700dbe37bd614cf39281dd5))
* adds Dockerfile improvements ([cededf7](https://github.com/steamban/one2n-sre-bootcamp/commit/cededf74f4fc8346fb3c4c3bfe3f769c6f5dfe42))
* adds proper digest usage in Dockerfile ([0315cb4](https://github.com/steamban/one2n-sre-bootcamp/commit/0315cb44f04d3b6366cbc3a5e20bb8c8c9307cac))
* **api:** adds pagination to student list api ([b09b2d9](https://github.com/steamban/one2n-sre-bootcamp/commit/b09b2d9ac7727e5a56df6ff760ecba9d554ae34d))
* **api:** updates repo func name and return type ([39fe8b2](https://github.com/steamban/one2n-sre-bootcamp/commit/39fe8b201ce89849c5657afd00a411cec1e62192))
* **api:** updates repo func name and return type ([c18c36f](https://github.com/steamban/one2n-sre-bootcamp/commit/c18c36fa9006e3adf88991b4b4a18835e62c148d))
* **ci:** remove unwanted build job ([fe86b54](https://github.com/steamban/one2n-sre-bootcamp/commit/fe86b548746d6ddcbdce9f2655943eba194d7678))
* **ci:** update failure threshold for hadolint ([705ff02](https://github.com/steamban/one2n-sre-bootcamp/commit/705ff02fda4651aaa0b9451015f8f3c32018b994))
* **ci:** update path filter to match actual code directories ([155a954](https://github.com/steamban/one2n-sre-bootcamp/commit/155a9540c5597aff5855583680b3a8e38194e5c2))
* **ci:** update release trigger tag ([e29ee39](https://github.com/steamban/one2n-sre-bootcamp/commit/e29ee39364104bb009a387c87055c80905df5b27))
* correct postgres healthcheck in docker compose ([58ab68d](https://github.com/steamban/one2n-sre-bootcamp/commit/58ab68d24a158438993daaaec582d60e13a44121))
* removes hardcoded healthcheck port for api ([0119922](https://github.com/steamban/one2n-sre-bootcamp/commit/011992270a381a769c91db8423fb68fe7fe2564a))
* resolves multiple issues in Makefile and docker compose file ([e1e273c](https://github.com/steamban/one2n-sre-bootcamp/commit/e1e273c1ae2040e620ef6baf5bdacd325c37191a))
* updates example env to match docker compose defaults ([e38f573](https://github.com/steamban/one2n-sre-bootcamp/commit/e38f5732e50c55e302a696b6d9afbfc91e6c7a48))
* use .env for database credentials in docker-compose ([e996b5a](https://github.com/steamban/one2n-sre-bootcamp/commit/e996b5a493b6042040eda0edf40dd7246926f663))
* use PORT env var for API port mapping ([4e4f6c2](https://github.com/steamban/one2n-sre-bootcamp/commit/4e4f6c2af84386055190f8253aec2aa32aa44fe4))
