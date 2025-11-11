<!--
SPDX-FileCopyrightText: 2022-present Intel Corporation
SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
Copyright 2019 free5GC.org

SPDX-License-Identifier: Apache-2.0
-->
[![Go Report Card](https://goreportcard.com/badge/github.com/omec-project/smf)](https://goreportcard.com/report/github.com/omec-project/smf)

# SMF

SMF is a session management function in 5G architecture and acts as the anchor
point to communicate with UPF (User Plane Function). SD-Core SMF supports
interaction with multiple UPFs. SMF supports handling dynamic slice
configuration, removal and modification.

## SMF Block diagram

![SMF Architecture](/docs/images/README-SMF.png)

SMF has configuration interface to handle slice configuration. Config service is
realised using the project - Config Service. SMF exports  metrics to prometheus.


## Repository Structure

Below is a high-level view of the repository and its main components:
```
.
├── callback                    # Contains HTTP callback APIs and routing definitions for NF subscriptions and notifications.
│   ├── api_default.go
│   ├── api_nf_subscribe_notify.go
│   └── routers.go
├── config                      # Includes SMF configuration files, such as default settings, UE routing information, and sample configurations for custom WebUI URLs. 
│   ├── smfcfg_with_custom_webui_url.yaml
│   ├── smfcfg.yaml
│   └── uerouting.yaml
├── consumer                    # Implements SMF’s interaction with other Network Functions (NFs), handling NF management, SM policy communication, and PDU session callbacks.
│   ├── nf_management.go
│   ├── nf_management_test.go
│   ├── nsmf_pdusession_callback.go
│   └── sm_policy.go
├── context                     # Core logic for session management and user plane context handling. Defines the SMF internal state, PFCP session contexts, IP allocation, GBR/MBR management, and UE routing logic.                      
│   ├── apply_action.go
│   ├── bp_manager.go
│   ├── context.go
│   ├── context_test.go
│   ├── datapath.go
│   ├── datapath_test.go
│   ├── db.go
│   ├── db_uptunnel.go
│   ├── gate.go
│   ├── gbr.go
│   ├── gsm_build.go
│   ├── gsm_handler.go
│   ├── interface.go
│   ├── ip_allocator.go
│   ├── ip_allocator_test.go
│   ├── mbr.go
│   ├── nf_profile.go
│   ├── ngap_build.go
│   ├── ngap_handler.go
│   ├── nodeid.go
│   ├── nodeid_test.go
│   ├── outer_header.go
│   ├── pco.go
│   ├── pfcp_rules.go
│   ├── pfcp_session_context.go
│   ├── qfi.go
│   ├── rrsp_flags.go
│   ├── sm_context.go
│   ├── snssai_dnn_smf_info.go
│   ├── snssai.go
│   ├── traffic_control_data.go
│   ├── ue_datapath.go
│   ├── ue_routing_manager.go
│   ├── ue_routing_manager_test.go
│   ├── upf.go
│   ├── user_plane_function_features.go
│   ├── user_plane_information.go
│   └── user_plane_information_test.go
├── design                      # Contains design documentation, including conceptual designs such as UPF selection strategy.  
│   └── upf-selection-design.md
├── Dockerfile
├── Dockerfile.fast
├── docs                        # Documentation resources (e.g., diagrams or reference images) for SMF documentation.
│   └── images
│       ├── README-SMF.png
│       └── README-SMF.png.license
├── eventexposure               # Provides APIs related to event exposure functions, enabling other NFs to subscribe or be notified about SMF-related events. 
│   ├── api_default.go
│   └── routers.go
├── factory                     # Responsible for loading and initializing configuration files and factory methods used during SMF startup.
│   ├── config.go
│   ├── config_test.go
│   └── factory.go
├── fsm                         # Contains the finite state machine (FSM) implementation that manages SMF session lifecycle transitions and state handling.
│   ├── handler.go
│   ├── txnFsm.go
│   └── utils.go
├── go.mod
├── go.mod.license
├── go.sum
├── go.sum.license
├── LICENSES
│   └── Apache-2.0.txt
├── logger                      # Centralized logging utilities used across the SMF for consistent log output and debugging.
│   └── logger.go
├── Makefile
├── metrics                     # Handles SMF telemetry and metrics export (e.g., via Kafka) for monitoring and observability.
│   ├── kafka.go
│   ├── kafka_test.go
│   └── telemetry.go
├── msgtypes                    # Defines service message types used internally for NF communication and signaling.
│   └── svcmsgtypes
│       └── svc_msgtypes.go
├── nfregistration              # Contains logic for SMF registration with the NRF (Network Repository Function), managing NF discovery and deregistration.
│   ├── nf_registration.go
│   └── nf_registration_test.go
├── NOTICE.txt
├── oam                         # Implements Operation, Administration, and Maintenance (OAM) APIs to retrieve session-related information for management tools.
│   ├── api_get_ue_pdu_session_info.go
│   └── routers.go
├── pdusession                  # Implements RESTful APIs for PDU Session Management, including SM Context creation, update, and release.
│   ├── api_individual_pdu_session_hsmf.go
│   ├── api_individual_sm_context.go
│   ├── api_pdu_sessions_collection.go
│   ├── api_sm_contexts_collection.go
│   ├── dummy_server.go
│   └── routers.go
├── pfcp                        # Implements the Packet Forwarding Control Protocol (PFCP) stack. Includes submodules for message handling, IEs, UDP transport, and communication with UPFs.
│   ├── adapter
│   │   └── adapter.go
│   ├── dispatcher.go
│   ├── handler
│   │   ├── handler.go
│   │   └── handler_test.go
│   ├── ies
│   │   ├── user_plane_function_features.go
│   │   └── user_plane_function_features_test.go
│   ├── message
│   │   ├── build.go
│   │   ├── build_test.go
│   │   ├── send.go
│   │   └── send_test.go
│   ├── udp
│   │   ├── event.go
│   │   ├── message.go
│   │   ├── pfcp_message_type.go
│   │   ├── pfcp_message_type_test.go
│   │   ├── transaction.go
│   │   ├── udp.go
│   │   └── udp_test.go
│   └── upf
│       └── upf.go
├── polling                     # Provides logic for SMF periodic configuration updates and NF status polling.
│   ├── nf_configuration.go
│   └── nf_configuration_test.go
├── producer                    # Contains SMF procedures that generate responses or events for other NFs, including PDU session creation, N1N2 signaling, and ULCL (Uplink Classifier) handling.
│   ├── callback.go
│   ├── datapath.go
│   ├── n1n2_data_handler.go
│   ├── oam.go
│   ├── pdu_session.go
│   ├── sm_pfcp_handling.go
│   ├── subscription_test.go
│   └── ulcl_procedure.go
├── qos                         # Manages QoS (Quality of Service) logic such as PCC rules, QoS flows, policies, and traffic control data
│   ├── charging_data.go
│   ├── condition_data.go
│   ├── pcc_rule.go
│   ├── qos_flow.go
│   ├── qos_flow_test.go
│   ├── qos_rule.go
│   ├── qos_rule_test.go
│   ├── qos_utility.go
│   ├── session_rule.go
│   ├── smPolicyData.go
│   └── traffic_control_data.go
├── README.md
├── service                     # Defines SMF initialization and startup logic, including server and API setup.
│   └── init.go
├── smferrors                   # Contains custom error definitions used across SMF modules.
│   └── errors.go
├── smf.go
├── Taskfile.yml
├── test-mirror.txt
├── transaction                 # Provides transaction management utilities for request handling and message correlation.
│   └── transaction.go
├── util                        # General-purpose utilities, including QoS conversion helpers.
│   └── qos_convert.go
├── VERSION
└── VERSION.license

33 directories, 132 files
```

## Configuration and Deployment

**Docker**

To build the container image:
```
task mod-start
task build
task docker-build-fast
```

**Kubernetes**

The standard deployment uses Helm charts from the Aether project. The version of the Chart can be found in the OnRamp repository in the `vars/main.yml` file.


## Quick Navigation

| Goal                                          | Path / Directory                                             | Description                                                           |
| --------------------------------------------- | ------------------------------------------------------------ | --------------------------------------------------------------------- |
| **Modify SMF configuration**                  | [`config/`](./config)                                        | Adjust default SMF settings, UE routing definitions, or WebUI URLs.   |
| **Understand PDU Session Management APIs**    | [`pdusession/`](./pdusession)                                | RESTful interfaces for PDU session lifecycle operations.              |
| **Review PFCP logic or UPF communication**    | [`pfcp/`](./pfcp)                                            | PFCP message encoding, dispatching, and UPF interaction.              |
| **Analyze SMF internal state or UE routing**  | [`context/`](./context)                                      | Core context management, PFCP sessions, and data path control.        |
| **Check NF registration and discovery logic** | [`nfregistration/`](./nfregistration)                        | Handles SMF registration, discovery, and heartbeat with NRF.          |
| **View SMF startup code**                     | [`service/init.go`](./service/init.go), [`smf.go`](./smf.go) | Entry points for initializing and running the SMF service.            |
| **Inspect QoS and PCC rule logic**            | [`qos/`](./qos)                                              | Policy and charging control, QoS flows, and session rules.            |
| **Explore metrics and telemetry exports**     | [`metrics/`](./metrics)                                      | Produces telemetry data for observability and performance monitoring. |
| **Find PDU session handling procedures**      | [`producer/`](./producer)                                    | Core SMF logic for creating and managing PDU sessions.                |
| **Check SMF design notes or diagrams**        | [`design/`](./design), [`docs/`](./docs)                     | Documentation and architecture design visuals.                        |
| **Build or containerize SMF**                 | [`Makefile`](./Makefile), [`Dockerfile`](./Dockerfile)       | Tools to compile, test, or deploy SMF as a containerized service.     |


## Supported Features
1. Supports PDU Session Establishment, Modification, Release
2. N2/X2 handover
3. End Marker Indication to UPF
4. PfcpSessionReport
5. N1N2MessageTransferFailureNotification handling Callback handling
6. Slice based UPF selection
7. UE address pool per Slice
8. PFCP heartbeat towards UPF
9. UE IP-Address allocation via UPF
10. QoS  call flows in SMF to handle PCC rules in Create Session Policy Response
and installing those rules in UPF & UE
11. High Availibilty and Cloud Native support(scale up/down number of instances
and subscriber store in Database)
12. UPF-Adapter for PFCP registration of multiple SMF instances with same
node-id to any UPF
13. Keep-alive support with respect to NRF
14. Transaction queueing for the same PDU session
15. SMF metrics available via metric-func to 5g Grafana dashboard
16. Static IP-address provision via configuration


## SMF supports wide range of error handling,
This includes some of the handling as listed below
1. UPF Reconnect if UPF restarts
2. PFCP Heartbeat handling towards UPF
3. PFCP Transaction timeout and not to wait forever
4. SBI message timeout handling and handling timeouts
5. Registration towards NRF with updated configuration
6. Retrying NRF registration if NRF is not available

## Upcoming features in SMF

1. Policy Notify from PCF for QoS update

Compliance of the 5G Network functions can be found at [5G Compliance](https://docs.sd-core.opennetworking.org/main/overview/3gpp-compliance-5g.html)

Design section for SMF is available at [SMF Design](https://docs.sd-core.opennetworking.org/main/design/design-smf.html)

## How to use SMF

Refer to the [SD-Core documentation](https://docs.sd-core.opennetworking.org/main/index.html)


## Reach out to us thorugh

1. #sdcore-dev channel in [ONF Community Slack](https://onf-community.slack.com/)
2. Raise Github issues

