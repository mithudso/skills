---
name: mac-egpu-compute
description: Instructions and deep research for configuring an external NVIDIA/AMD GPU for compute tasks on Apple Silicon Macs using the TinyGPU driver.
---

# Mac eGPU Compute via TinyGPU

## Overview
Historically, Apple Silicon Macs (M1/M2/M3/M4) dropped all support for external GPUs (eGPUs). However, in April 2026, the **tinygrad** team (Tiny Corp) released the `TinyGPU` driver, which Apple officially signed. 

This driver allows Apple Silicon to communicate with external AMD (RDNA3+) and NVIDIA (Ampere+) GPUs over Thunderbolt/USB4 for **compute-only tasks** (e.g., running LLMs, AI inference). It does *not* provide display output or graphics acceleration.

## Requirements
*   **Mac:** Apple Silicon running macOS 13.0 or later.
*   **GPU:** AMD RDNA3+ or NVIDIA Ampere+ inside a USB4/Thunderbolt enclosure.
*   **Docker:** Required for NVIDIA configurations (not required for AMD).

## Installation Procedure

### Step 1: Install the TinyGPU Driver
Run the official setup script from tinygrad:
```bash
curl -fsSL https://raw.githubusercontent.com/tinygrad/tinygrad/master/extra/setup_tinygpu_osx.sh | sh
```

### Step 2: Enable Driver Extension
If not prompted automatically, the user must go to **System Settings > General > Login Items & Extensions > Driver Extensions** and toggle "TinyGPU" to ON.

### Step 3: Install Docker Desktop (NVIDIA Only)
NVIDIA GPUs require the use of Docker Desktop, as the NVIDIA compiler (NVCC) tooling relies on a Linux container environment.

### Step 4: Run NVCC Setup (NVIDIA Only)
Run the NVIDIA compiler setup script:
```bash
curl -fsSL https://raw.githubusercontent.com/tinygrad/tinygrad/master/extra/setup_nvcc_osx.sh | sh
```

### Step 5: Run Compute Tasks
Once the driver is active and the container is bridged, you can run LLM inference tools (like Ollama or vLLM) inside the Docker container, accessing the eGPU hardware.

## Troubleshooting
* **SIP:** System Integrity Protection (SIP) does *not* need to be disabled because the driver is notarized by Apple.
* **Bandwidth:** Thunderbolt limits bandwidth to ~40 Gbps, which may cause slightly slower model load times compared to native PCIe, though token generation speed remains high.
