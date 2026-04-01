#!/bin/bash
lvm help > /tmp/lvm_help.txt 2>&1 || man -k lvm > /tmp/lvm_help.txt 2>&1 || echo "LVM - Logical Volume Manager" > /tmp/lvm_help.txt
dd if=/dev/zero of=/tmp/disk1.img bs=1M count=100 2>/dev/null
dd if=/dev/zero of=/tmp/disk2.img bs=1M count=100 2>/dev/null
cat > /tmp/lvm_concepts.txt << LVMEOF
LVM Concepts:
- PV (Physical Volume): 物理卷，对应物理磁盘或分区
- VG (Volume Group): 卷组，由一个或多个PV组成
- LV (Logical Volume): 逻辑卷，从VG中分配空间
流程: PV -> VG -> LV -> 文件系统
LVMEOF
