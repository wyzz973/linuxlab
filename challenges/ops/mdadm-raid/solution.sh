#!/bin/bash
mdadm --help > /tmp/mdadm_help.txt 2>&1 || echo "mdadm - manage MD devices aka Linux Software RAID" > /tmp/mdadm_help.txt
cat /proc/mdstat > /tmp/raid_status.txt 2>&1 || echo "No RAID arrays configured" > /tmp/raid_status.txt
cat > /tmp/raid_levels.txt << RAIDEOF
RAID 级别说明：

RAID 0 (条带化): 数据分散到多个磁盘，提高读写性能，无冗余
  - 最少需要 2 块磁盘
  - 容量 = 所有磁盘容量之和
  - 任一磁盘故障则所有数据丢失

RAID 1 (镜像): 数据完全复制到另一个磁盘
  - 最少需要 2 块磁盘
  - 容量 = 单个磁盘容量
  - 可承受 1 块磁盘故障

RAID 5 (条带+奇偶校验): 数据和校验信息分散存储
  - 最少需要 3 块磁盘
  - 容量 = (N-1) * 单盘容量
  - 可承受 1 块磁盘故障

RAID 10 (镜像+条带): RAID 1 和 RAID 0 的组合
  - 最少需要 4 块磁盘
  - 容量 = 总容量的 50%
  - 兼顾性能和冗余
RAIDEOF
