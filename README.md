# Avalanche Simplified

## Steps to run

### 1. Change settings

| Environment | Config File                               |
|-------------|-------------------------------------------|
| Docker      | `src/worker/config/prod.yml`              |
| Testing     | `src/worker/config/testing.yml`           |
| E2E-Testing | `src/worker/functional-tests/default.yml` |

### 2. Steps to tests E2E

#### 2.0. Prerequisite

Run `echo 50000 > /proc/sys/kernel/keys/maxkeys`, otherwise some containers would not be able to start.

#### 2.1. Setup container

Run `make build-docker`

Run `run-200-worker-docker`

#### 2.2. Create terminal to the `client_node` container and run

Run `make test-e2e`

## Explain

### 1. Nodes exchange their IP addresses by Broadcasting a message

Refer to the `src/worker/jobs/self_introduction.go` file.

### 2. After collecting enough neighbor IP addresses, the nodes start to exchange messages to reach consensus

Refer to the `src/worker/jobs/fetch_transaction.go` file.
