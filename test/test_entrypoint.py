"""Run with: python3 -m unittest discover -s test -p 'test_entrypoint.py'."""

import os
from pathlib import Path
import subprocess
import tempfile
import unittest


class EntrypointTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.action = self.root / "action"
        for size in ("dev", "test"):
            (self.action / size).mkdir(parents=True)
        self.log = self.root / "calls"
        rc = self.root / "bashrc"
        rc.touch()
        for tool in ("terraform", "pulumi"):
            stub = self.root / tool
            stub.write_text(
                '#!/bin/sh\n'
                'printf "%s:%s:%s\\n" "${0##*/}" "$PWD" "$*" >> "$CALL_LOG"\n'
            )
            stub.chmod(0o755)
        source = (Path(__file__).resolve().parents[1] / "entrypoint.sh").read_text()
        # Substitute container-only paths; execute the actual validation and dispatch.
        source = source.replace("/root/.bashrc", str(rc))
        source = source.replace("/action/", str(self.action) + "/")
        source = source.replace("/pulumi/bin/pulumi", str(self.root / "pulumi"))
        self.script = self.root / "entrypoint.sh"
        self.script.write_text(source)

    def run_entrypoint(self, **inputs):
        if self.log.exists():
            self.log.unlink()
        env = {key: value for key, value in os.environ.items()
               if not key.startswith("INPUT_")}
        env.update({
            "PATH": str(self.root) + os.pathsep + os.defpath,
            "CALL_LOG": str(self.log),
            "INPUT_CLUSTER_SIZE": "dev",
            "INPUT_ACTION_TYPE": "create",
            "INPUT_IAC_TYPE": "terraform",
        })
        env.update(inputs)
        result = subprocess.run(
            ["/bin/sh", str(self.script)], cwd=self.root, env=env,
            capture_output=True, text=True, check=False,
        )
        calls = self.log.read_text().splitlines() if self.log.exists() else []
        return result, calls

    def test_invalid_inputs_stop_before_iac(self):
        values = {
            "CLUSTER_SIZE": ("", "../test", "../../../github/workspace/attacker_iac", "DEV"),
            "ACTION_TYPE": ("", "delete", "CREATE", "create ", "create; id", "create\nid"),
            "IAC_TYPE": ("", "other", "TERRAFORM", "terraform ", "pulumi; id", "pulumi\nid"),
        }
        for name, invalid_values in values.items():
            for value in invalid_values:
                with self.subTest(input=name, value=value):
                    result, calls = self.run_entrypoint(**{"INPUT_" + name: value})
                    self.assertNotEqual(result.returncode, 0)
                    self.assertIn("Invalid " + name, result.stderr)
                    self.assertEqual(calls, [])

    def test_valid_inputs_dispatch_from_selected_directory(self):
        for size in ("dev", "test"):
            for tool in ("terraform", "pulumi"):
                for action in ("create", "destroy"):
                    with self.subTest(size=size, tool=tool, action=action):
                        result, calls = self.run_entrypoint(
                            INPUT_CLUSTER_SIZE=size, INPUT_IAC_TYPE=tool,
                            INPUT_ACTION_TYPE=action,
                        )
                        self.assertNotIn("Invalid", result.stderr)
                        self.assertTrue(calls)
                        prefix = f"{tool}:{self.action / size}:"
                        self.assertTrue(all(call.startswith(prefix) for call in calls))
                        command = calls[-1].removeprefix(prefix)
                        expected = {
                            ("terraform", "create"): "apply tfplan",
                            ("terraform", "destroy"): "destroy -force",
                            ("pulumi", "create"): "up --yes",
                            ("pulumi", "destroy"): "destroy -force",
                        }[(tool, action)]
                        self.assertEqual(command, expected)

    def test_missing_directory_stops_before_iac(self):
        (self.action / "dev").rmdir()
        for tool in ("terraform", "pulumi"):
            with self.subTest(tool=tool):
                result, calls = self.run_entrypoint(INPUT_IAC_TYPE=tool)
                self.assertNotEqual(result.returncode, 0)
                self.assertEqual(calls, [])


if __name__ == "__main__":
    unittest.main()
