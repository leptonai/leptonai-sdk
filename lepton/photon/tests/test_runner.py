import os
import tempfile

# Set cache dir to a temp dir before importing anything from lepton
tmpdir = tempfile.mkdtemp()
os.environ["LEPTON_CACHE_DIR"] = tmpdir

import json
from textwrap import dedent
import sys
import unittest
import zipfile

from loguru import logger
import numpy as np
import requests
import torch

import lepton
from lepton.photon import RunnerPhoton as Runner


from utils import random_name, photon_run_server


class CustomRunner(Runner):
    def init(self):
        self.nn = torch.nn.Linear(1, 1)

    @Runner.handler("some_path")
    def run(self, x: float) -> float:
        return self.nn(torch.tensor(x).reshape(1, 1)).item()


class CustomRunnerWithCustomDeps(Runner):
    requirement_dependency = ["torch"]

    def init(self):
        self.nn = torch.nn.Linear(1, 1)

    @Runner.handler("some_path")
    def run(self, x: float) -> float:
        return self.nn(torch.tensor(x).reshape(1, 1)).item()


class TestRunner(unittest.TestCase):
    def test_run(self):
        name = random_name()
        runner = CustomRunner(name=name)
        x = 2.0
        y1 = runner.run(x)

        xtensor = torch.tensor(x).reshape(1, 1)
        y2 = runner.nn(xtensor).item()
        self.assertEqual(y1, y2)

    def test_save_load(self):
        name = random_name()
        runner = CustomRunner(name=name)
        x = 2.0
        y1 = runner.run(x)

        path = runner.save()

        runner = lepton.photon.load(path)
        y2 = runner.run(x)
        self.assertEqual(y1, y2)

    def test_run_server(self):
        # pytest imports test files as top-level module which becomes
        # unavailable in server process
        if "PYTEST_CURRENT_TEST" in os.environ:
            import cloudpickle

            cloudpickle.register_pickle_by_value(sys.modules[__name__])

        name = random_name()
        runner = CustomRunner(name=name)
        path = runner.save()

        proc, port = photon_run_server(path=path)

        x = 2.0
        res = requests.post(
            f"http://localhost:{port}/some_path",
            json={"x": x},
        )
        proc.kill()
        self.assertEqual(res.status_code, 200)

    def test_runner_cli(self):
        with tempfile.NamedTemporaryFile(suffix=".py") as f:
            f.write(
                dedent(
                    """
from lepton.photon.runner import RunnerPhoton as Runner, handler


class Counter(Runner):
    def init(self):
        self.counter = 0

    @handler("add")
    def add(self, x: int) -> int:
        self.counter += x
        return self.counter

    @handler("sub")
    def sub(self, x: int) -> int:
        self.counter -= x
        return self.counter
"""
                ).encode("utf-8")
            )
            f.flush()
            proc, port = photon_run_server(name="counter", model=f"py:{f.name}:Counter")
            res = requests.post(
                f"http://127.0.0.1:{port}/add",
                json={"x": 1},
            )
            self.assertEqual(res.status_code, 200)
            self.assertEqual(res.json(), 1)

            res = requests.post(
                f"http://127.0.0.1:{port}/add",
                json={"x": 1},
            )
            self.assertEqual(res.status_code, 200)
            self.assertEqual(res.json(), 2)

            res = requests.post(
                f"http://127.0.0.1:{port}/sub",
                json={"x": 2},
            )
            self.assertEqual(res.status_code, 200)
            self.assertEqual(res.json(), 0)
            proc.kill()

    def test_photon_file_metadata(self):
        name = random_name()
        runner = CustomRunner(name=name)
        path = runner.save()
        with zipfile.ZipFile(path, "r") as photon_file:
            with photon_file.open("metadata.json") as metadata_file:
                metadata = json.load(metadata_file)
        self.assertEqual(metadata["name"], name)
        self.assertEqual(metadata["model"], "CustomRunner")
        self.assertTrue("image" in metadata)
        self.assertTrue("args" in metadata)
        self.assertTrue("openapi" in metadata)
        self.assertTrue("/some_path" in metadata["openapi"]["paths"])
        self.assertGreater(len(metadata.get("requirement_dependency")), 0)

    def test_custom_requirement_dependency(self):
        name = random_name()
        runner = CustomRunnerWithCustomDeps(name=name)
        path = runner.save()
        with zipfile.ZipFile(path, "r") as photon_file:
            with photon_file.open("metadata.json") as metadata_file:
                metadata = json.load(metadata_file)
        self.assertEqual(
            metadata["requirement_dependency"],
            CustomRunnerWithCustomDeps.requirement_dependency,
        )

    def test_metrics(self):
        # pytest imports test files as top-level module which becomes
        # unavailable in server process
        if "PYTEST_CURRENT_TEST" in os.environ:
            import cloudpickle

            cloudpickle.register_pickle_by_value(sys.modules[__name__])

        name = random_name()
        runner = CustomRunner(name=name)
        path = runner.save()

        proc, port = photon_run_server(path=path)

        for x in range(5):
            res = requests.post(
                f"http://127.0.0.1:{port}/some_path",
                json={"x": float(x)},
            )
            self.assertEqual(res.status_code, 200)
        res = requests.get(f"http://127.0.0.1:{port}/metrics")
        self.assertEqual(res.status_code, 200)
        self.assertRegex(
            res.text, r'http_request_duration_seconds_count{handler="/some_path"}'
        )
        proc.kill()

    def test_vec_db_examples(self):
        name = random_name()

        vec_db_path = os.path.join(
            os.path.dirname(lepton.__file__), "examples", "vec_db.py"
        )
        with tempfile.NamedTemporaryFile(suffix=".py", dir=tmpdir) as f, open(
            vec_db_path, "rb"
        ) as vec_db_file:
            f.write(vec_db_file.read())
            f.flush()
            proc, port = photon_run_server(name="vec_db", model=f"py:{f.name}:VecDB")

        dim = 2
        name = "two"

        # create collection
        res = requests.post(
            f"http://127.0.0.1:{port}/create_collection",
            json={"name": name, "dim": dim},
        )
        self.assertEqual(res.status_code, 200)

        # list collections
        # TODO: this should be get, not post
        res = requests.post(f"http://127.0.0.1:{port}/list_collections", json={})
        self.assertEqual(res.status_code, 200)
        self.assertEqual(res.json(), [[name, dim]])

        # create second collection, list collections, remove it and list collections again
        name2 = random_name()
        res = requests.post(
            f"http://127.0.0.1:{port}/create_collection",
            json={"name": name2, "dim": dim},
        )
        self.assertEqual(res.status_code, 200)
        res = requests.post(f"http://127.0.0.1:{port}/list_collections", json={})
        self.assertEqual(res.status_code, 200)
        self.assertTrue([name2, dim] in res.json())
        res = requests.post(
            f"http://127.0.0.1:{port}/remove_collection", json={"name": name2}
        )
        self.assertEqual(res.status_code, 200)
        res = requests.post(f"http://127.0.0.1:{port}/list_collections", json={})
        self.assertEqual(res.status_code, 200)
        self.assertEqual(res.json(), [[name, dim]])

        # insert
        count = 10
        embeddings = []
        for i in range(count):
            vector = np.random.rand(dim).tolist()
            text = f"text_{i}"
            doc_id = f"doc_id_{i}"
            embeddings.append({"doc_id": doc_id, "text": text, "vector": vector})
        res = requests.post(
            f"http://127.0.0.1:{port}/add",
            json={"name": name, "embeddings": embeddings},
        )
        self.assertEqual(res.status_code, 200)
        res = requests.post(f"http://127.0.0.1:{port}/count", json={"name": name})
        self.assertEqual(res.status_code, 200)
        self.assertEqual(res.json(), count)

        # search
        k = 3
        res = requests.post(
            f"http://127.0.0.1:{port}/search",
            json={"name": name, "vector": embeddings[0]["vector"], "k": k},
        )
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(res.json()), k)
        self.assertEqual(res.json()[0]["doc_id"], embeddings[0]["doc_id"])

        # delete
        res = requests.post(
            f"http://127.0.0.1:{port}/delete",
            json={"name": name, "doc_ids": [embeddings[0]["doc_id"]]},
        )
        self.assertEqual(res.status_code, 200)
        res = requests.post(f"http://127.0.0.1:{port}/count", json={"name": name})
        self.assertEqual(res.status_code, 200)
        self.assertEqual(res.json(), count - 1)

        proc.kill()


if __name__ == "__main__":
    unittest.main()
