"""customer-db — the PII store. One benign skill, two dangerous ones."""

import os

from agentfield import Agent

app = Agent(
    node_id="customer-db",
    agentfield_server=os.getenv("AGENTFIELD_URL", "http://localhost:8080"),
    tags=["pii"],
    enable_did=True,
    vc_enabled=True,
)

_DB = {
    "CUST-A": {"name": "Alice", "email": "alice@acme.com", "ssn": "***-**-1111"},
    "CUST-B": {"name": "Bob", "email": "bob@acme.com", "ssn": "***-**-2222"},
}


@app.skill(tags=["pii"])
def get_record(customer_id: str, limit: int = 1) -> dict:
    return {"customer_id": customer_id, **_DB.get(customer_id, {})}


@app.skill(tags=["pii"])
def export_all() -> dict:
    return {"customers": _DB}  # the bulk-leak weapon


@app.skill(tags=["pii"])
def delete_record(customer_id: str) -> dict:
    return {"deleted": customer_id}


if __name__ == "__main__":
    app.run(auto_port=True)
