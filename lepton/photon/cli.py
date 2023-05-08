import os
import sys

from rich.console import Console
from rich.table import Table
import click
from .base import find_all_photons, find_photon, remove_photon
from . import api

console = Console(highlight=False)


@click.group()
def photon():
    pass


@photon.command()
@click.option("--name", "-n", help="Name of the Photon")
@click.option("--model", "-m", help="Model spec")
def create(name, model):
    console.print(f"Creating Photon: [green]{name}[/green]")
    try:
        photon = api.create(name=name, model=model)
    except Exception as e:
        console.print(f"Failed to create Photon:\n{e}")
        sys.exit(1)
    try:
        api.save(photon)
    except Exception as e:
        console.print(f'Failed to save Photon: "{e}"')
        sys.exit(1)
    console.print(f"Photon [green]{name}[/green] created")


@photon.command()
@click.option("--name", "-n", help="Name of the Photon")
def remove(name):
    if find_photon(name) is None:
        console.print(f'Photon "{name}" [red]does not exist[/]')
        sys.exit(1)
    remove_photon(name)
    console.print(f'Photon "{name}" [green]removed[/]')


@photon.command()
@click.option(
    "--remote_url", "-r", help="Remote URL of the Lepton Server", default=None
)
def list(remote_url):
    if remote_url is not None:
        # TODO: Add Creation Time and other metadata
        photons = api.list_remote(remote_url)
        records = [
            (photon["name"], photon["model"], photon["id"]) for photon in photons
        ]
    else:
        records = find_all_photons()
        records = [
            (name, model, id_) for id_, name, model, path, creation_time in records
        ]

    table = Table(title="Photons", show_lines=True)
    table.add_column("Name")
    table.add_column("Model")
    table.add_column("ID")

    records_by_name = {}
    for name, model, id_ in records:
        records_by_name.setdefault(name, []).append((model, id_))
    for name, sub_records in records_by_name.items():
        model_table = Table(show_header=False, box=None)
        id_table = Table(show_header=False, box=None)
        for model, id_ in sub_records:
            model_table.add_row(model)
            id_table.add_row(id_)
        table.add_row(name, model_table, id_table)
    console.print(table)


@photon.command()
@click.option("--name", "-n", help="Name of the Photon", default=None)
@click.option("--model", "-m", help="Model Spec", default=None)
@click.option("--file", "-f", "path", help="Path to .photon file", default=None)
@click.option("--port", "-p", help="Port to run on", default=8080)
@click.pass_context
def run(ctx, name, model, path, port):
    if name is None and path is None:
        console.print("Must specify either --name or --path")
        sys.exit(1)
    if path is None:
        path = find_photon(name)
    if path is None or not os.path.exists(path):
        name_or_path = name if name is not None else path
        console.print(f'Photon "{name_or_path}" [red]does not exist[/]')
        if name and model:
            ctx.invoke(create, name=name, model=model)
            path = find_photon(name)
        else:
            sys.exit(1)
    photon = api.load(path)
    photon.launch(port=port)


@photon.command()
@click.option("--name", "-n", help="Name of the Photon")
@click.option("--remote_url", "-r", help="Remote URL of the Lepton Server")
def push(name, remote_url):
    path = find_photon(name)
    if path is None or not os.path.exists(path):
        console.print(f'Photon "{name}" [red]does not exist[/]')
        sys.exit(1)
    api.push(path, remote_url)
    console.print(f'Photon "{name}" [green]pushed[/]')


def add_command(click_group):
    click_group.add_command(photon)
