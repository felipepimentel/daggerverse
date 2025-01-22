package main

import (
  "dagger.io/dagger"
)

dagger.#Plan & {
  actions: {
    "test-deploy": {
      n8n: dagger.#N8N
      deploy: n8n.#Deploy & {
        input: {
          domain:       "test.example.com"
          subdomain:    "n8n-test"
          n8nVersion:   "0.234.0"
          region:       "syd1"
          dropletSize:  "s-1vcpu-1gb"
          doToken:      dagger.#Secret
        }
      }
    }
  }
}