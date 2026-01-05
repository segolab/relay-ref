# Infrastructure

This directory contains Infrastructure-as-Code (IaC) definitions for
deploying the Relay reference implementation.

Infrastructure is organized by IaC tool and orchestration platform.
Each lane is independent and optional.

## Structure

terraform/
  Kubernetes and Nomad definitions using Terraform

pulumi-go/ (planned)
  Kubernetes and Nomad definitions using Pulumi (Go)

pulumi-cs/ (planned)
  Kubernetes and Nomad definitions using Pulumi (C#)

## Baseline

The initial baseline focuses on:
- Terraform
- Kubernetes
- Local development (e.g. kind)

Additional lanes are introduced incrementally and must not be required
to build, test, or run the baseline.

## Lane independence

Unused infrastructure lanes can be ignored or removed without breaking
active lanes or automation.
