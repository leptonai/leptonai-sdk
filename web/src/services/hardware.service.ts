import { Hardware } from "@lepton-dashboard/interfaces/hardware";
import { Injectable } from "injection-js";

@Injectable()
export class HardwareService {
  hardwareShapes: Hardware = {
    "cpu.small": {
      Description:
        "General purpose small shape with 1 CPUs, 4GB of RAM and 16GB of ephemeral storage",
      DisplayName: "cpu.small",
      Selectable: true,
      Resource: {
        CPU: 1,
        Memory: 4096,
        EphemeralStorageInGB: 16,
      },
    },
    "cpu.medium": {
      Description:
        "General purpose medium shape with 2 CPUs, 8GB of RAM and 32GB of ephemeral storage",
      DisplayName: "cpu.medium",
      Selectable: true,
      Resource: {
        CPU: 2,
        Memory: 8192,
        EphemeralStorageInGB: 32,
      },
    },
    "cpu.large": {
      Description:
        "General purpose large shape with 4 CPUs, 16GB of RAM and 64GB of ephemeral storage",
      DisplayName: "cpu.large",
      Selectable: true,
      Resource: {
        CPU: 4,
        Memory: 16384,
        EphemeralStorageInGB: 64,
      },
    },
    "gpu.t4": {
      Description:
        "Accelerated computing shape with 1 16GB T4 GPU, 4 CPUs, 16GB of RAM and 100GB of ephemeral storage",
      DisplayName: "gpu.t4",
      Selectable: true,
      Resource: {
        CPU: 4,
        Memory: 16384,
        EphemeralStorageInGB: 100,
        AcceleratorType: "Tesla-T4",
        AcceleratorNum: 1,
      },
    },
    "gpu.a10": {
      Description:
        "Accelerated computing shape with 1 24GB A10 GPU, 8 CPUs, 32GB of RAM and 400GB of ephemeral storage",
      DisplayName: "gpu.a10",
      Selectable: true,
      Resource: {
        CPU: 8,
        Memory: 32768,
        EphemeralStorageInGB: 400,
        AcceleratorType: "NVIDIA-A10",
        AcceleratorNum: 1,
      },
    },
    "ac1.a10": {
      Description:
        "Accelerated computing shape with 1 24GB A10 GPU, 8 CPUs, 32GB of RAM and 400GB of ephemeral storage",
      DisplayName: "ac1.a10",
      Selectable: false,
      Resource: {
        CPU: 8,
        Memory: 32768,
        EphemeralStorageInGB: 400,
        AcceleratorType: "NVIDIA-A10",
        AcceleratorNum: 1,
      },
    },
  };
  isGPUInstance(shape?: string): boolean {
    return this.hardwareShapes && shape
      ? !!this.hardwareShapes[shape]?.Resource?.AcceleratorType
      : false;
  }

  shapes = Object.entries(this.hardwareShapes)
    .filter(([_, shape]) => shape.Selectable)
    .map(([shape]) => shape);
}
