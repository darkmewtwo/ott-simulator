Run from `api/` with Python 3.12 or newer:

```powershell
py -3.12 -m venv .venv
.venv/Scripts/python.exe -m pip install -r requirements-dev.txt
.venv/Scripts/python.exe -m pytest -q
```

These unit tests mock repositories and use temporary upload directories. They
do not need PostgreSQL, Docker, or a running streaming service.

Run the simulator suite from `simulator/`:

```powershell
go test ./... -cover
```

Simulator tests use local HTTP test servers and require no running API.
