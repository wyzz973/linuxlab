#!/usr/bin/env python3
"""Drive the LinuxLab React TUI under a PTY to verify the Go handoff chain.

Scenario: boot -> modules -> first challenge -> detail -> reveal hint (h)
-> Enter (start challenge) -> shell handoff -> `ls > /tmp/result.txt; exit`
-> result screen. Prints captured output and exits non-zero on missing markers.
"""
import os
import pty
import select
import sys
import time

CMD = sys.argv[1] if len(sys.argv) > 1 else "npm run tui:react"
# 进入挑战 shell 后执行的命令（默认适用于 linux-basics 类题目；
# compose 挑战可传 "docker compose up -d; exit"）
SHELL_CMD = os.environ.get("SMOKE_SHELL_CMD", "ls > /tmp/result.txt; exit")
COLS, ROWS = 110, 30
MARKERS = [
    b"LinuxLab",
    b"\xe8\xae\xad\xe7\xbb\x83\xe6\x8e\xa7\xe5\x88\xb6\xe5\x8f\xb0",  # 训练控制台
    b"\xe5\xbc\x80\xe5\xa7\x8b\xe7\xbb\x83\xe4\xb9\xa0",              # 开始练习
]
DETAIL_MARKERS = [
    b"\xe6\x8f\x90\xe7\xa4\xba",                                      # 提示
    b"\xe6\x8c\x89 h \xe6\x9f\xa5\xe7\x9c\x8b\xe7\xac\xac\xe4\xb8\x80\xe6\x9d\xa1\xe6\x8f\x90\xe7\xa4\xba",  # 按 h 查看第一条提示
    b"Enter \xe5\xbc\x80\xe5\xa7\x8b\xe6\x8c\x91\xe6\x88\x98",         # Enter 开始挑战
]
HINT_MARKERS = [
    b"\xe5\xb7\xb2\xe7\x94\xa8 1/",                                    # 已用 1/
    b"\xe5\xbd\xb1\xe5\x93\x8d\xe5\xbe\x97\xe5\x88\x86",              # 影响得分
]
RESULT_MARKERS = [
    b"\xe6\xa3\x80\xe6\xb5\x8b\xe7\xbb\x93\xe6\x9e\x9c",              # 检测结果
    b"\xe9\x80\x9a\xe8\xbf\x87",                                      # 通过
]
DOCKER_PS1 = b"[linuxlab]"

pid, fd = pty.fork()
if pid == 0:
    env = dict(os.environ)
    env.update({
        "COLUMNS": str(COLS),
        "LINES": str(ROWS),
        "LINUXLAB_BIN": "./linuxlab",
        "TERM": "xterm-256color",
        "NO_COLOR": "0",
    })
    os.execvpe("/bin/sh", ["/bin/sh", "-c", CMD], env)

output = bytearray()
start = time.time()
timeout = 60


def pump(seconds):
    end = time.time() + seconds
    while time.time() < end:
        r, _, _ = select.select([fd], [], [], 0.2)
        if r:
            try:
                chunk = os.read(fd, 4096)
            except OSError:
                return False
            if not chunk:
                return False
            output.extend(chunk)
    return True


def send(data):
    if isinstance(data, str):
        data = data.encode("utf-8")
    os.write(fd, data)


def wait_for(markers, seconds, label):
    deadline = time.time() + seconds
    found = []
    while time.time() < deadline:
        if any(m in output for m in markers):
            found = [m for m in markers if m in output]
            return found
        if not pump(0.3):
            break
    return found


ok = True
step = "boot"


def check(step_name, found, expected):
    global ok, step
    step = step_name
    if not found:
        ok = False
        print(f"[FAIL] {step_name}: markers missing; expected one of {expected}")
    else:
        print(f"[ OK ] {step_name}: {found}")


try:
    found = wait_for(MARKERS, 25, "boot")
    check("boot menu", found, MARKERS)
    send("2")  # 开始练习 -> modules
    pump(1.0)
    send("\r")  # enter first module (Linux 基础命令)
    pump(1.0)
    send("\r")  # enter first challenge detail
    pump(1.5)
    found = wait_for(DETAIL_MARKERS, 10, "detail")
    check("detail screen", found, DETAIL_MARKERS)
    send("h")  # reveal hint
    pump(1.0)
    found = wait_for(HINT_MARKERS, 8, "hint reveal")
    check("hint reveal", found, HINT_MARKERS)
    send("\r")  # Enter: start challenge -> Go handoff
    found = wait_for([DOCKER_PS1], 30, "shell handoff")
    check("shell handoff", found, [DOCKER_PS1])
    pump(1.0)
    send(SHELL_CMD + "\r")
    found = wait_for(RESULT_MARKERS, 30, "result screen")
    check("result screen", found, RESULT_MARKERS)
    pump(1.0)
    send("q")  # back from result
    pump(1.0)
    send("q")
    pump(0.5)
    send("\x03")  # safety Ctrl+C
except Exception as exc:  # noqa: BLE001
    ok = False
    print(f"[FAIL] exception during {step}: {exc}")
finally:
    try:
        os.kill(pid, 15)
    except ProcessLookupError:
        pass
    try:
        os.close(fd)
    except OSError:
        pass

print(f"\n=== elapsed {time.time() - start:.1f}s, total bytes {len(output)} ===")
with open("/tmp/react-tui-smoke.log", "wb") as f:
    f.write(bytes(output))
if ok:
    print("SMOKE TEST PASSED")
else:
    print("SMOKE TEST FAILED at step:", step)
sys.exit(0 if ok else 1)
