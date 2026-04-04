#!/bin/bash
ansible --version > /tmp/ansible_version.txt 2>&1 || echo "Ansible version: conceptual exercise - ansible not installed in this container" > /tmp/ansible_version.txt

cat > /tmp/inventory.ini << INVEOF
[webservers]
web1.example.com ansible_host=192.168.1.10
web2.example.com ansible_host=192.168.1.11

[dbservers]
db1.example.com ansible_host=192.168.1.20

[all:vars]
ansible_user=deploy
ansible_python_interpreter=/usr/bin/python3
INVEOF

cat > /tmp/setup.yml << YMLEOF
---
- name: Setup web servers
  hosts: webservers
  become: yes
  tasks:
    - name: Update apt cache
      apt:
        update_cache: yes

    - name: Install nginx
      apt:
        name: nginx
        state: present

    - name: Start nginx
      service:
        name: nginx
        state: started
        enabled: yes
YMLEOF
