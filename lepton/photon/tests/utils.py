import atexit
from contextlib import closing
import os
import random
import socket
import subprocess
import string
import time


def find_free_port():
    # ref: https://stackoverflow.com/a/45690594
    with closing(socket.socket(socket.AF_INET, socket.SOCK_STREAM)) as s:
        s.bind(('', 0))
        s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        return s.getsockname()[1]


def random_name():
    return "".join(random.choice(string.ascii_lowercase) for _ in range(5))


def photon_run_server(name=None, path=None, model=None, port=None):
    if name is None and path is None:
        raise ValueError("Either name or path must be specified")
    if name is not None and path is not None:
        raise ValueError("Only one of name or path can be specified")
    if port is None:
        port = find_free_port()
    cmd = [
        "lepton",
        "photon",
        "run",
    ]
    if name:
        cmd += ["-n", name]
    if path:
        cmd += ["-f", path]
    cmd += ["--port", str(port)]
    if model is not None:
        cmd += ["-m", model]

    proc = subprocess.Popen(
        cmd,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    atexit.register(proc.kill)
    lines = []
    for line in proc.stderr:
        line = line.decode("utf-8")
        lines.append(line)
        if "running" in line.lower():
            break
        time.sleep(0.1)
    else:
        # "running" never showed up in the output, which means the
        # server failed to start
        proc.kill()
        stdout = proc.stdout.read().decode("utf-8")
        stderr = os.linesep.join(lines)
        raise RuntimeError(
            f"Photon server failed to start:\nstdout:\n{stdout}\nstderr:\n{stderr}"
        )
    return proc, port
