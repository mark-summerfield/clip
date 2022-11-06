#!/usr/bin/env python3
# Copyright © 2022 Mark Summerfield. All rights reserved.
# License: GPLv3

import os
import subprocess


def main():
    os.chdir(os.path.dirname(__file__))
    total = ok = 0
    for cmd, expected in (
            (['eg/eg1/eg1.go'], 'eg1.txt'),
            (['eg/eg1/eg1.go', '-h'], 'eg1.hlp'),
            (['eg/eg2/eg2.go'], 'eg2.txt'),
            (['eg/eg2/eg2.go', '-h'], 'eg2.hlp'),
            (['eg/subcommands/subcommands.go', '-h'], 'subcommands.hlp'),
            (['eg/subcommands/subcommands.go', '-h', 'c'],
             'subcommands-c.hlp'),
            (['eg/subcommands/subcommands.go', '-h', 'l'],
             'subcommands-l.hlp'),
            (['eg/subcommands/subcommands.go', '-h', 'f'],
             'subcommands-f.hlp'),
            ):
        total += 1
        ok += test(cmd, expected)
    if total != ok:
        print(f'FAIL {ok}/{total}')
    else:
        print(f'All {total} OK')


def test(cmd, expected):
    cmd = ['go', 'run'] + cmd
    reply = subprocess.run(cmd, capture_output=True, text=True,
                           encoding='utf-8')
    cmd = ' '.join(cmd)
    if reply.returncode:
        print(f'FAIL bad return code "{cmd}": {reply.returncode}')
        return 0
    if reply.stderr:
        print(f'FAIL stderr "{cmd}": {reply.stderr}')
        return 0
    try:
        filename = f'rdata/{expected}'
        with open(filename, 'rt', encoding='utf-8') as file:
            text = file.read().strip()
        if reply.stdout.strip() != text:
            print(f'FAIL stdout != {filename}')
            return 0
        return 1
    except FileNotFoundError:
        print(f'FAIL missing: {filename}')
        return 0


if __name__ == '__main__':
    main()
