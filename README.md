# Avalanche Simplified
This is a README file for the Avalanche Simplified project. It provides instructions for running the program and some explanation about the program's functionality.

## Steps to Run

### 1. Change Settings

To change the program's settings, you need to edit the appropriate configuration file. The following table lists the configuration files for each environment:

| Environment | Config File                        |
|-------------|------------------------------------|
| Docker      | `src/worker/config/prod.yml`       |
| Testing     | `src/worker/config/testing.yml`    |
| E2E-Testing | `src/worker/e2e-tests/default.yml` |

### 2. Steps to Test E2E

Before testing E2E, you need to complete the following prerequisite step:

#### 2.0. Prerequisite

Run `echo 50000 > /proc/sys/kernel/keys/maxkeys` to allow containers to start.

To test E2E, follow these steps:

#### 2.1. Setup Container

IMPORTANT: Run `make down-200-worker-docker` between each run. The program assumes that nodes are spawned incrementally from 1 to N.


1. Run `make build-docker`
2. Run `run-200-worker-docker`

#### 2.2. Create Terminal to the `client_node` Container and Run

1. Create a terminal to the `client_node` container
2. Run `make test-e2e`

## Explanation

The Avalanche Simplified project is based on a network of nodes that exchange messages to reach consensus. The following are some details about the program's functionality:

### 1. Nodes Update Their Neighbors' Liveness by Sending a REST Request to Each Other

The `src/worker/jobs/scan_nodes.go` file contains the code that handles this feature.

### 2. Nodes Exchange Messages to Reach Consensus

The `src/worker/jobs/fetch_transaction.go` file contains the code that allows nodes to exchange messages and reach consensus.
